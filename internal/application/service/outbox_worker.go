package service

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"time"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type OutboxWorker struct {
	repo     DB.OutboxRepositoryInterface
	users    DB.UserInterface
	payments DB.PaymentRepositoryInterface
	tickets  DB.TicketRepositoryInterface
	issuer   service_contract.TicketDocumentIssuer
	email    service_contract.EmailServiceProviderInterface
	log      domainlogger.Logger
}

func NewOutboxWorker(repo DB.OutboxRepositoryInterface, users DB.UserInterface, payments DB.PaymentRepositoryInterface, tickets DB.TicketRepositoryInterface, issuer service_contract.TicketDocumentIssuer, email service_contract.EmailServiceProviderInterface, log domainlogger.Logger) *OutboxWorker {
	return &OutboxWorker{repo: repo, users: users, payments: payments, tickets: tickets, issuer: issuer, email: email, log: log}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		if err := w.runBatch(ctx); err != nil && ctx.Err() == nil {
			w.log.Error(ctx, "outbox batch failed", domainlogger.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *OutboxWorker) runBatch(ctx context.Context) error {
	events, err := w.repo.Claim(ctx, 10, time.Now())
	if err != nil {
		return err
	}
	for i := range events {
		eventCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		err := w.process(eventCtx, &events[i])
		cancel()
		if err == nil {
			if markErr := w.repo.MarkProcessed(ctx, events[i].ID, time.Now()); markErr != nil {
				return markErr
			}
			continue
		}
		delay := time.Duration(1<<min(events[i].Attempts, 9)) * time.Second
		if markErr := w.repo.MarkFailed(ctx, events[i].ID, err.Error(), time.Now().Add(delay)); markErr != nil {
			return markErr
		}
		w.log.Warn(ctx, "outbox event processing failed", domainlogger.Int64("event_id", events[i].ID), domainlogger.String("event_type", events[i].EventType), domainlogger.Error(err))
	}
	return nil
}

func (w *OutboxWorker) process(ctx context.Context, event *entity.OutboxEvent) error {
	var payload struct {
		PaymentID int64 `json:"payment_id"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil || payload.PaymentID <= 0 {
		return fmt.Errorf("invalid outbox payload")
	}
	switch event.EventType {
	case entity.OutboxTicketPDFGenerate:
		return w.issuer.EnsurePDFForPayment(ctx, payload.PaymentID)
	case entity.OutboxEmailPaymentSucceeded:
		return w.sendPaymentEmail(ctx, payload.PaymentID, false)
	case entity.OutboxEmailPaymentRefunded:
		return w.sendPaymentEmail(ctx, payload.PaymentID, true)
	default:
		return fmt.Errorf("unsupported outbox event type %q", event.EventType)
	}
}

func (w *OutboxWorker) sendPaymentEmail(ctx context.Context, paymentID int64, refunded bool) error {
	payment, err := w.payments.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	user, err := w.users.FindUserByID(ctx, payment.UserID)
	if err != nil {
		return err
	}
	if user.Email == nil || *user.Email == "" || !user.IsEmailVerified {
		w.log.Info(ctx, "email notification skipped because user has no verified email", domainlogger.Int64("user_id", user.ID), domainlogger.Int64("payment_id", paymentID))
		return nil
	}
	subject := "Payment successful - Cinema Booking System"
	statusText := "Your payment was successful and your ticket is ready."
	if refunded {
		subject = "Payment refunded - Cinema Booking System"
		statusText = "Your payment was refunded successfully."
	}
	body := fmt.Sprintf("%s\nPayment ID: %d\nAmount: %d", statusText, payment.ID, payment.Amount)
	htmlBody := fmt.Sprintf("<p>%s</p><p>Payment ID: <strong>%d</strong></p><p>Amount: <strong>%d</strong></p>", html.EscapeString(statusText), payment.ID, payment.Amount)
	_, err = w.email.SendEmail(ctx, dto.SendEmailRequest{To: []string{*user.Email}, Subject: subject, TextBody: body, HTMLBody: htmlBody})
	return err
}

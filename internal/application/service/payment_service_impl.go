package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/broker"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type paymentService struct {
	repo            DB.PaymentRepositoryInterface
	reservationRepo DB.ReservationRepositoryInterface
	gateway         broker.PaymentGateway
	log             domainlogger.Logger
}

func NewPaymentService(repo DB.PaymentRepositoryInterface, reservationRepo DB.ReservationRepositoryInterface, gateway broker.PaymentGateway, log domainlogger.Logger) service_contract.PaymentServiceInterface {
	return &paymentService{repo: repo, reservationRepo: reservationRepo, gateway: gateway, log: log}
}

func (s *paymentService) Create(ctx context.Context, reservationID, userID int64, idempotencyKey string) (*dto.PaymentResponse, bool, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if reservationID <= 0 || userID <= 0 || idempotencyKey == "" || utf8.RuneCountInString(idempotencyKey) > 100 {
		return nil, false, exception.ErrInvalidInput
	}
	reservation, err := s.reservationRepo.GetReservation(ctx, reservationID)
	if err != nil {
		return nil, false, err
	}
	if reservation.UserID != userID {
		return nil, false, exception.ErrNotAccessible
	}
	authority, _, err := s.gateway.Create(ctx, reservationID, reservation.TotalAmount)
	if err != nil {
		return nil, false, exception.Wrap(exception.ErrInternal, "mock payment initialization failed", err)
	}
	payment, created, err := s.repo.CreatePayment(ctx, &entity.Payment{
		ReservationID:  reservationID,
		UserID:         userID,
		Provider:       enum.PaymentProviderMock,
		Status:         enum.PaymentPending,
		Authority:      authority,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return nil, false, err
	}
	s.log.Info(ctx, "Payment initialized successfully",
		domainlogger.Int64("payment_id", payment.ID),
		domainlogger.Int64("reservation_id", reservationID),
		domainlogger.Int64("user_id", userID),
		domainlogger.Bool("created", created),
	)
	return dto.ToPaymentResponse(payment), created, nil
}

func (s *paymentService) Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.PaymentResponse, error) {
	if id <= 0 || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	payment, err := s.repo.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && payment.UserID != userID {
		return nil, exception.ErrNotAccessible
	}
	return dto.ToPaymentResponse(payment), nil
}

func (s *paymentService) ListForReservation(ctx context.Context, reservationID, userID int64, allowAny bool) ([]*dto.PaymentResponse, error) {
	if reservationID <= 0 || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	reservation, err := s.reservationRepo.GetReservation(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	if !allowAny && reservation.UserID != userID {
		return nil, exception.ErrNotAccessible
	}
	payments, err := s.repo.ListReservationPayments(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.PaymentResponse, 0, len(payments))
	for i := range payments {
		result = append(result, dto.ToPaymentResponse(&payments[i]))
	}
	return result, nil
}

func (s *paymentService) Verify(ctx context.Context, id, userID int64, allowAny bool, req *dto.VerifyMockPaymentRequest) (*dto.PaymentResponse, error) {
	if id <= 0 || userID <= 0 || req == nil || req.Success == nil {
		return nil, exception.ErrInvalidInput
	}
	payment, err := s.repo.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowAny && payment.UserID != userID {
		return nil, exception.ErrNotAccessible
	}
	if payment.Status != enum.PaymentPending {
		return dto.ToPaymentResponse(payment), nil
	}

	result, err := s.gateway.Verify(ctx, payment.Authority, *req.Success)
	if err != nil {
		return nil, exception.Wrap(exception.ErrInternal, "mock payment verification failed", err)
	}
	if result == nil || (result.Succeeded && strings.TrimSpace(result.ReferenceID) == "") {
		return nil, exception.Wrap(exception.ErrInternal, "mock payment gateway returned an invalid result", errors.New("invalid gateway result"))
	}
	failureReason := strings.TrimSpace(req.FailureReason)
	if failureReason == "" {
		failureReason = result.FailureReason
	}
	finalized, err := s.repo.FinalizePayment(ctx, id, userID, allowAny, result.Succeeded, result.ReferenceID, failureReason, time.Now())
	if err != nil {
		return nil, err
	}
	s.log.Info(ctx, "Payment finalized successfully",
		domainlogger.Int64("payment_id", finalized.ID),
		domainlogger.Int64("reservation_id", finalized.ReservationID),
		domainlogger.Int64("user_id", userID),
		domainlogger.String("status", string(finalized.Status)),
	)
	return dto.ToPaymentResponse(finalized), nil
}

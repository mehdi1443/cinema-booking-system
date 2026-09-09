package postgres_repo

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestOperationsOutboxAndPaginationIntegration(t *testing.T) {
	dsn := os.Getenv("PAYMENT_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("PAYMENT_INTEGRATION_DSN is not set")
	}
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()
	wrapped := &database.PostgresDatabase{DB: tx}
	ctx := context.Background()
	userID, _, reservationID := createPaymentFixture(t, tx)
	paymentRepo := NewPaymentRepository(wrapped)
	payment, created, err := paymentRepo.CreatePayment(ctx, &entity.Payment{ReservationID: reservationID, UserID: userID, Authority: "ops-authority", IdempotencyKey: "ops-key"})
	if err != nil || !created {
		t.Fatalf("create payment: created=%v err=%v", created, err)
	}
	paidAt := time.Now()
	if _, err := paymentRepo.FinalizePayment(ctx, payment.ID, userID, false, true, "ops-reference", "", paidAt); err != nil {
		t.Fatalf("finalize payment: %v", err)
	}

	payments, paymentTotal, err := paymentRepo.ListReservationPayments(ctx, reservationID, 1, 0)
	if err != nil || len(payments) != 1 || paymentTotal != 1 {
		t.Fatalf("payment pagination invalid: len=%d total=%d err=%v", len(payments), paymentTotal, err)
	}
	bookingRepo := NewBookingRepository(wrapped)
	reservations, reservationTotal, err := bookingRepo.ListUserReservations(ctx, userID, 1, 0)
	if err != nil || len(reservations) != 1 || reservationTotal != 1 {
		t.Fatalf("reservation pagination invalid: len=%d total=%d err=%v", len(reservations), reservationTotal, err)
	}
	ticketRepo := NewTicketRepository(wrapped)
	tickets, ticketTotal, err := ticketRepo.ListUserTickets(ctx, userID, 1, 0)
	if err != nil || len(tickets) != 1 || ticketTotal != 1 {
		t.Fatalf("ticket pagination invalid: len=%d total=%d err=%v", len(tickets), ticketTotal, err)
	}

	outboxRepo := NewOutboxRepository(wrapped)
	events, err := outboxRepo.Claim(ctx, 10, time.Now().Add(time.Second))
	if err != nil || len(events) != 2 {
		t.Fatalf("transactional outbox invalid: len=%d err=%v", len(events), err)
	}
	for _, event := range events {
		if err := outboxRepo.MarkProcessed(ctx, event.ID, time.Now()); err != nil {
			t.Fatal(err)
		}
	}

	var cinemaID, showtimeID int64
	if err := tx.Raw(`SELECT h.cinema_id, r.showtime_id FROM reservations r JOIN showtimes s ON s.id=r.showtime_id JOIN halls h ON h.id=s.hall_id WHERE r.id=?`, reservationID).Row().Scan(&cinemaID, &showtimeID); err != nil {
		t.Fatal(err)
	}
	if err := tx.Model(&entity.Showtime{}).Where("id = ?", showtimeID).Updates(map[string]any{"starts_at": paidAt.Add(-2 * time.Hour), "ends_at": paidAt.Add(-time.Hour)}).Error; err != nil {
		t.Fatal(err)
	}
	ops := NewOperationsRepository(wrapped)
	report, err := ops.SalesReport(ctx, cinemaID, paidAt.Add(-time.Minute), paidAt.Add(time.Minute))
	if err != nil || report.NetAmount != payment.Amount || report.SuccessfulPayments != 1 {
		t.Fatalf("sales report invalid: %#v err=%v", report, err)
	}
	settlement, err := ops.CreateSettlement(ctx, cinemaID, userID, paidAt.Add(-time.Minute), paidAt.Add(time.Minute))
	if err != nil || settlement.GrossAmount != payment.Amount {
		t.Fatalf("create settlement: %#v err=%v", settlement, err)
	}
	if _, err := ops.CreateSettlement(ctx, cinemaID, userID, paidAt.Add(-time.Minute), paidAt.Add(time.Minute)); !errors.Is(err, exception.ErrNoSettlementPayments) {
		t.Fatalf("payment settled twice: %v", err)
	}
	paid, err := ops.MarkSettlementPaid(ctx, settlement.ID, userID, time.Now())
	if err != nil || paid.Status != "PAID" || paid.PaidAt == nil {
		t.Fatalf("mark paid: %#v err=%v", paid, err)
	}
	settlements, total, err := ops.ListSettlements(ctx, cinemaID, 1, 0)
	if err != nil || len(settlements) != 1 || total != 1 {
		t.Fatalf("settlement pagination invalid: len=%d total=%d err=%v", len(settlements), total, err)
	}

	entry := &entity.AuditLog{ActorUserID: &userID, Method: "POST", Path: "/v1/test", StatusCode: 201}
	if err := ops.RecordAudit(ctx, entry); err != nil {
		t.Fatal(err)
	}
	audits, auditTotal, err := ops.ListAudit(ctx, DB.AuditFilter{ActorUserID: &userID, Limit: 1})
	if err != nil || len(audits) != 1 || auditTotal < 1 {
		t.Fatalf("audit persistence invalid: len=%d total=%d err=%v", len(audits), auditTotal, err)
	}
}

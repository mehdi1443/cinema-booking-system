package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPaymentLifecycleIntegration(t *testing.T) {
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
		t.Fatalf("begin test transaction: %v", tx.Error)
	}
	defer tx.Rollback()

	ctx := context.Background()
	userID, otherUserID, reservationID := createPaymentFixture(t, tx)
	repo := NewPaymentRepository(&database.PostgresDatabase{DB: tx})

	first, created, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: reservationID, UserID: userID, Authority: "mock_first_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-1",
	})
	if err != nil || !created {
		t.Fatalf("create first payment: created=%v err=%v", created, err)
	}
	if first.Amount != 250_000 || first.Status != enum.PaymentPending || first.Provider != enum.PaymentProviderMock {
		t.Fatalf("server payment values are invalid: %#v", first)
	}

	replayed, created, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: reservationID, UserID: userID, Authority: "must_be_ignored", IdempotencyKey: "checkout-1",
	})
	if err != nil || created || replayed.ID != first.ID {
		t.Fatalf("idempotent replay failed: payment=%#v created=%v err=%v", replayed, created, err)
	}
	if _, err := repo.FinalizePayment(ctx, first.ID, otherUserID, false, true, "unauthorized-ref", "", time.Now()); !errors.Is(err, exception.ErrNotAccessible) {
		t.Fatalf("another user finalized payment: %v", err)
	}

	failed, err := repo.FinalizePayment(ctx, first.ID, userID, false, false, "", "DECLINED_FOR_TEST", time.Now())
	if err != nil || failed.Status != enum.PaymentFailed || failed.FailureReason == nil {
		t.Fatalf("failed verification was not persisted: payment=%#v err=%v", failed, err)
	}
	assertReservationAndSeatStatus(t, tx, reservationID, enum.ReservationPendingPayment, enum.ShowtimeSeatHeld)

	second, created, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: reservationID, UserID: userID, Authority: "mock_second_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-2",
	})
	if err != nil || !created {
		t.Fatalf("create retry payment: created=%v err=%v", created, err)
	}
	succeeded, err := repo.FinalizePayment(ctx, second.ID, userID, false, true, "mock-reference-success", "", time.Now())
	if err != nil || succeeded.Status != enum.PaymentSucceeded || succeeded.ReferenceID == nil || succeeded.PaidAt == nil {
		t.Fatalf("successful verification failed: payment=%#v err=%v", succeeded, err)
	}
	assertReservationAndSeatStatus(t, tx, reservationID, enum.ReservationPaid, enum.ShowtimeSeatSold)
	var ticket entity.Ticket
	if err := tx.Where("payment_id = ?", second.ID).First(&ticket).Error; err != nil {
		t.Fatalf("successful payment did not issue a ticket: %v", err)
	}
	if ticket.Status != enum.TicketActive || ticket.UserID != userID || ticket.ReservationID != reservationID || ticket.Code == "" {
		t.Fatalf("issued ticket is invalid: %#v", ticket)
	}
	ticketRepo := NewTicketRepository(&database.PostgresDatabase{DB: tx})
	if _, err := ticketRepo.ConsumeTicket(ctx, ticket.Code, otherUserID, false, time.Now()); !errors.Is(err, exception.ErrNotAccessible) {
		t.Fatalf("non-owner cinema user validated ticket: %v", err)
	}
	usedTicket, err := ticketRepo.ConsumeTicket(ctx, ticket.Code, userID, false, time.Now())
	if err != nil || usedTicket.Status != enum.TicketUsed || usedTicket.UsedAt == nil || usedTicket.UsedByUserID == nil {
		t.Fatalf("cinema owner could not validate ticket: ticket=%#v err=%v", usedTicket, err)
	}
	if _, err := ticketRepo.ConsumeTicket(ctx, ticket.Code, userID, false, time.Now()); !errors.Is(err, exception.ErrTicketAlreadyUsed) {
		t.Fatalf("ticket was not one-time use: %v", err)
	}

	replayedFinal, err := repo.FinalizePayment(ctx, second.ID, userID, false, true, "different-reference", "", time.Now())
	if err != nil || replayedFinal.ReferenceID == nil || *replayedFinal.ReferenceID != "mock-reference-success" {
		t.Fatalf("finalization replay was not idempotent: payment=%#v err=%v", replayedFinal, err)
	}

	expiringUserID, _, expiringReservationID := createPaymentFixture(t, tx)
	expiringPayment, _, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: expiringReservationID, UserID: expiringUserID, Authority: "mock_expiring_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-expiring",
	})
	if err != nil {
		t.Fatalf("create expiring payment: %v", err)
	}
	past := time.Now().Add(-time.Minute)
	if err := tx.Model(&entity.Reservation{}).Where("id = ?", expiringReservationID).Update("expires_at", past).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Model(&entity.ShowtimeSeat{}).Where("reservation_id = ?", expiringReservationID).Update("held_until", past).Error; err != nil {
		t.Fatal(err)
	}
	bookingRepo := NewBookingRepository(&database.PostgresDatabase{DB: tx})
	if err := bookingRepo.ExpireReservations(ctx, time.Now()); err != nil {
		t.Fatalf("expire reservation: %v", err)
	}
	expiredPayment, err := repo.GetPayment(ctx, expiringPayment.ID)
	if err != nil || expiredPayment.Status != enum.PaymentFailed || expiredPayment.FailureReason == nil || *expiredPayment.FailureReason != "RESERVATION_EXPIRED" {
		t.Fatalf("reservation expiry did not fail payment: payment=%#v err=%v", expiredPayment, err)
	}
	assertReservationAndSeatStatus(t, tx, expiringReservationID, enum.ReservationExpired, enum.ShowtimeSeatAvailable)

	refundUserID, refundOtherUserID, refundReservationID := createPaymentFixture(t, tx)
	refundablePayment, _, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: refundReservationID, UserID: refundUserID, Authority: "mock_refund_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-refund",
	})
	if err != nil {
		t.Fatalf("create refundable payment: %v", err)
	}
	if _, err := repo.FinalizePayment(ctx, refundablePayment.ID, refundUserID, false, true, "mock-reference-refund", "", time.Now()); err != nil {
		t.Fatalf("finalize refundable payment: %v", err)
	}
	if _, err := repo.RefundPayment(ctx, refundablePayment.ID, refundOtherUserID, false, "unauthorized-refund", time.Now(), 0); !errors.Is(err, exception.ErrNotAccessible) {
		t.Fatalf("another user refunded payment: %v", err)
	}
	refunded, err := repo.RefundPayment(ctx, refundablePayment.ID, refundUserID, false, "mock-refund-reference", time.Now(), 0)
	if err != nil || refunded.Status != enum.PaymentRefunded || refunded.RefundedAt == nil || refunded.RefundReferenceID == nil {
		t.Fatalf("refund failed: payment=%#v err=%v", refunded, err)
	}
	assertReservationAndSeatStatus(t, tx, refundReservationID, enum.ReservationCancelled, enum.ShowtimeSeatAvailable)
	var cancelledTicket entity.Ticket
	if err := tx.Where("payment_id = ?", refundablePayment.ID).First(&cancelledTicket).Error; err != nil || cancelledTicket.Status != enum.TicketCancelled {
		t.Fatalf("refund did not cancel ticket: ticket=%#v err=%v", cancelledTicket, err)
	}
	replayedRefund, err := repo.RefundPayment(ctx, refundablePayment.ID, refundUserID, false, "different-refund-reference", time.Now(), 0)
	if err != nil || replayedRefund.RefundReferenceID == nil || *replayedRefund.RefundReferenceID != "mock-refund-reference" {
		t.Fatalf("refund replay was not idempotent: payment=%#v err=%v", replayedRefund, err)
	}

	lateUserID, _, lateReservationID := createPaymentFixture(t, tx)
	latePayment, _, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: lateReservationID, UserID: lateUserID, Authority: "mock_late_refund_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-late-refund",
	})
	if err != nil {
		t.Fatalf("create late refund payment: %v", err)
	}
	if _, err := repo.FinalizePayment(ctx, latePayment.ID, lateUserID, false, true, "mock-reference-late-refund", "", time.Now()); err != nil {
		t.Fatalf("finalize late refund payment: %v", err)
	}
	if _, err := repo.RefundPayment(ctx, latePayment.ID, lateUserID, false, "mock-late-refund", time.Now(), 2*time.Hour); !errors.Is(err, exception.ErrRefundNotAllowed) {
		t.Fatalf("late refund was accepted: %v", err)
	}
	assertReservationAndSeatStatus(t, tx, lateReservationID, enum.ReservationPaid, enum.ShowtimeSeatSold)

	cancelledUserID, _, cancelledReservationID := createPaymentFixture(t, tx)
	cancelledPayment, _, err := repo.CreatePayment(ctx, &entity.Payment{
		ReservationID: cancelledReservationID, UserID: cancelledUserID, Authority: "mock_showtime_cancel_" + fmt.Sprint(time.Now().UnixNano()), IdempotencyKey: "checkout-showtime-cancel",
	})
	if err != nil {
		t.Fatalf("create showtime cancellation payment: %v", err)
	}
	if _, err := repo.FinalizePayment(ctx, cancelledPayment.ID, cancelledUserID, false, true, "mock-reference-showtime-cancel", "", time.Now()); err != nil {
		t.Fatalf("finalize showtime cancellation payment: %v", err)
	}
	var cancelledReservation entity.Reservation
	if err := tx.First(&cancelledReservation, cancelledReservationID).Error; err != nil {
		t.Fatal(err)
	}
	showtime, err := bookingRepo.GetShowtime(ctx, cancelledReservation.ShowtimeID)
	if err != nil {
		t.Fatalf("get showtime for cancellation: %v", err)
	}
	showtime.Status = enum.ShowtimeStatusCancelled
	if _, err := bookingRepo.UpdateShowtime(ctx, showtime); err != nil {
		t.Fatalf("cancel paid showtime: %v", err)
	}
	automaticRefund, err := repo.GetPayment(ctx, cancelledPayment.ID)
	if err != nil || automaticRefund.Status != enum.PaymentRefunded || automaticRefund.RefundedAt == nil || automaticRefund.RefundReferenceID == nil {
		t.Fatalf("showtime cancellation did not refund payment: payment=%#v err=%v", automaticRefund, err)
	}
	assertReservationAndSeatStatus(t, tx, cancelledReservationID, enum.ReservationCancelled, enum.ShowtimeSeatAvailable)
	var automaticallyCancelledTicket entity.Ticket
	if err := tx.Where("payment_id = ?", cancelledPayment.ID).First(&automaticallyCancelledTicket).Error; err != nil || automaticallyCancelledTicket.Status != enum.TicketCancelled {
		t.Fatalf("showtime cancellation did not cancel ticket: ticket=%#v err=%v", automaticallyCancelledTicket, err)
	}
}

func createPaymentFixture(t *testing.T, tx *gorm.DB) (int64, int64, int64) {
	t.Helper()
	stamp := time.Now().UnixNano()
	var userID, otherUserID, provinceID, cityID, cinemaID, movieID, hallID, seatID, showtimeID, reservationID, showtimeSeatID int64
	queries := []struct {
		query string
		args  []any
		dest  *int64
	}{
		{"INSERT INTO users (phone) VALUES (?) RETURNING id", []any{fmt.Sprintf("09%09d", stamp%1_000_000_000)}, &userID},
		{"INSERT INTO users (phone) VALUES (?) RETURNING id", []any{fmt.Sprintf("08%09d", stamp%1_000_000_000)}, &otherUserID},
		{"INSERT INTO provinces (name) VALUES (?) RETURNING id", []any{fmt.Sprintf("payment-province-%d", stamp)}, &provinceID},
	}
	for _, item := range queries {
		if err := tx.Raw(item.query, item.args...).Scan(item.dest).Error; err != nil {
			t.Fatalf("create payment fixture: %v", err)
		}
	}
	if err := tx.Raw("INSERT INTO cities (province_id, name) VALUES (?, ?) RETURNING id", provinceID, fmt.Sprintf("payment-city-%d", stamp)).Scan(&cityID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Raw("INSERT INTO cinemas (city_id, owner_user_id, name, address) VALUES (?, ?, ?, ?) RETURNING id", cityID, userID, fmt.Sprintf("payment-cinema-%d", stamp), "test").Scan(&cinemaID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Raw("INSERT INTO movies (title, duration_minutes, status) VALUES (?, 90, 'NOW_SHOWING') RETURNING id", fmt.Sprintf("payment-movie-%d", stamp)).Scan(&movieID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Raw("INSERT INTO halls (cinema_id, name, status) VALUES (?, ?, 'ACTIVE') RETURNING id", cinemaID, fmt.Sprintf("hall-%d", stamp)).Scan(&hallID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Raw("INSERT INTO seats (hall_id, row_label, number, seat_type, status) VALUES (?, 'A', 1, 'STANDARD', 'ACTIVE') RETURNING id", hallID).Scan(&seatID).Error; err != nil {
		t.Fatal(err)
	}
	startsAt := time.Now().Add(30 * time.Minute)
	if err := tx.Raw("INSERT INTO showtimes (hall_id, movie_id, starts_at, ends_at, status) VALUES (?, ?, ?, ?, 'PUBLISHED') RETURNING id", hallID, movieID, startsAt, startsAt.Add(90*time.Minute)).Scan(&showtimeID).Error; err != nil {
		t.Fatal(err)
	}
	expiresAt := time.Now().Add(20 * time.Minute)
	if err := tx.Raw("INSERT INTO reservations (user_id, showtime_id, status, expires_at, total_amount) VALUES (?, ?, 'PENDING_PAYMENT', ?, 250000) RETURNING id", userID, showtimeID, expiresAt).Scan(&reservationID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Raw("INSERT INTO showtime_seats (showtime_id, seat_id, row_label, seat_number, seat_type, price, status, reservation_id, held_until) VALUES (?, ?, 'A', 1, 'STANDARD', 250000, 'HELD', ?, ?) RETURNING id", showtimeID, seatID, reservationID, expiresAt).Scan(&showtimeSeatID).Error; err != nil {
		t.Fatal(err)
	}
	if err := tx.Exec("INSERT INTO reservation_items (reservation_id, showtime_seat_id, unit_price) VALUES (?, ?, 250000)", reservationID, showtimeSeatID).Error; err != nil {
		t.Fatal(err)
	}
	return userID, otherUserID, reservationID
}

func assertReservationAndSeatStatus(t *testing.T, tx *gorm.DB, reservationID int64, reservationStatus enum.ReservationStatus, seatStatus enum.ShowtimeSeatStatus) {
	t.Helper()
	var actualReservationStatus enum.ReservationStatus
	if err := tx.Raw("SELECT status FROM reservations WHERE id = ?", reservationID).Scan(&actualReservationStatus).Error; err != nil {
		t.Fatal(err)
	}
	var actualSeatStatus enum.ShowtimeSeatStatus
	if err := tx.Raw("SELECT ss.status FROM showtime_seats ss JOIN reservation_items ri ON ri.showtime_seat_id = ss.id WHERE ri.reservation_id = ?", reservationID).Scan(&actualSeatStatus).Error; err != nil {
		t.Fatal(err)
	}
	if actualReservationStatus != reservationStatus || actualSeatStatus != seatStatus {
		t.Fatalf("unexpected statuses: reservation=%s seat=%s", actualReservationStatus, actualSeatStatus)
	}
}

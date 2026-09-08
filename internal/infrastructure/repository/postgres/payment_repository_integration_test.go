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
	startsAt := time.Now().Add(2 * time.Hour)
	if err := tx.Raw("INSERT INTO showtimes (hall_id, movie_id, starts_at, ends_at, status) VALUES (?, ?, ?, ?, 'PUBLISHED') RETURNING id", hallID, movieID, startsAt, startsAt.Add(90*time.Minute)).Scan(&showtimeID).Error; err != nil {
		t.Fatal(err)
	}
	expiresAt := time.Now().Add(time.Hour)
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

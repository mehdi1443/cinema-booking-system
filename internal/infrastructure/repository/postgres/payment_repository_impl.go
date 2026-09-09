package postgres_repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentRepository struct {
	db database.Database
}

func NewPaymentRepository(db database.Database) DB.PaymentRepositoryInterface {
	return &paymentRepository{db: db}
}

func handlePaymentDBError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			switch pgErr.ConstraintName {
			case "uq_payments_user_idempotency":
				return exception.Wrap(exception.ErrIdempotencyConflict, exception.ErrIdempotencyConflict.Message, err)
			case "uq_payments_one_active_per_reservation":
				return exception.Wrap(exception.ErrPaymentInProgress, exception.ErrPaymentInProgress.Message, err)
			default:
				return exception.Wrap(exception.ErrDuplicateKey, exception.ErrDuplicateKey.Message, err)
			}
		case "23503":
			if strings.Contains(pgErr.ConstraintName, "reservation") {
				return exception.Wrap(exception.ErrReservationNotFound, exception.ErrReservationNotFound.Message, err)
			}
			if strings.Contains(pgErr.ConstraintName, "user") {
				return exception.Wrap(exception.ErrUserNotFound, exception.ErrUserNotFound.Message, err)
			}
		}
	}
	return exception.Wrap(exception.ErrDatabase, exception.ErrDatabase.Message, err)
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, bool, error) {
	if payment == nil {
		return nil, false, exception.ErrInvalidInput
	}
	created := false
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		now := time.Now()
		if err := expireReservations(tx, now); err != nil {
			return handlePaymentDBError(err)
		}

		var existing entity.Payment
		err := tx.Where("user_id = ? AND idempotency_key = ?", payment.UserID, payment.IdempotencyKey).First(&existing).Error
		if err == nil {
			if existing.ReservationID != payment.ReservationID {
				return exception.ErrIdempotencyConflict
			}
			*payment = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return handlePaymentDBError(err)
		}

		// Read the relation first, then lock showtime and reservation in that
		// order. Showtime cancellation uses the same order, preventing a
		// payment/finalization deadlock under concurrent checkout and admin edits.
		var reservationRef entity.Reservation
		if err := tx.Select("id, user_id, showtime_id").First(&reservationRef, payment.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&showtime, reservationRef.ShowtimeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handlePaymentDBError(err)
		}
		var reservation entity.Reservation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, payment.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		if reservation.UserID != payment.UserID {
			return exception.ErrNotAccessible
		}
		if showtime.Status != enum.ShowtimeStatusPublished || !showtime.StartsAt.After(now) {
			return exception.ErrShowtimeNotBookable
		}
		// The reservation lock serializes concurrent checkout requests. Rechecking the
		// key here makes simultaneous retries idempotent as well as sequential ones.
		err = tx.Where("user_id = ? AND idempotency_key = ?", payment.UserID, payment.IdempotencyKey).First(&existing).Error
		if err == nil {
			if existing.ReservationID != payment.ReservationID {
				return exception.ErrIdempotencyConflict
			}
			*payment = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return handlePaymentDBError(err)
		}
		if reservation.Status == enum.ReservationExpired || !reservation.ExpiresAt.After(now) {
			return exception.ErrReservationExpired
		}
		if reservation.Status != enum.ReservationPendingPayment {
			return exception.ErrReservationNotPayable
		}

		var active entity.Payment
		err = tx.Where("reservation_id = ? AND status IN ?", reservation.ID, []enum.PaymentStatus{enum.PaymentPending, enum.PaymentSucceeded}).First(&active).Error
		if err == nil {
			return exception.ErrPaymentInProgress
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return handlePaymentDBError(err)
		}

		var itemCount, heldCount int64
		if err := tx.Model(&entity.ReservationItem{}).Where("reservation_id = ?", reservation.ID).Count(&itemCount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		if err := tx.Model(&entity.ShowtimeSeat{}).Where("reservation_id = ? AND status = ? AND held_until > ?", reservation.ID, enum.ShowtimeSeatHeld, now).Count(&heldCount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		if itemCount == 0 || itemCount != heldCount {
			return exception.ErrReservationNotPayable
		}

		payment.Amount = reservation.TotalAmount
		payment.Status = enum.PaymentPending
		payment.Provider = enum.PaymentProviderMock
		if err := tx.Create(payment).Error; err != nil {
			return handlePaymentDBError(err)
		}
		created = true
		return nil
	})
	if err != nil {
		if errors.Is(err, exception.ErrIdempotencyConflict) {
			var existing entity.Payment
			lookupErr := r.db.GetDB().WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", payment.UserID, payment.IdempotencyKey).First(&existing).Error
			if lookupErr == nil && existing.ReservationID == payment.ReservationID {
				return &existing, false, nil
			}
		}
		return nil, false, err
	}
	return payment, created, nil
}

func (r *paymentRepository) GetPayment(ctx context.Context, id int64) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.GetDB().WithContext(ctx).First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrPaymentNotFound
		}
		return nil, handlePaymentDBError(err)
	}
	return &payment, nil
}

func (r *paymentRepository) ListReservationPayments(ctx context.Context, reservationID int64) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := r.db.GetDB().WithContext(ctx).Where("reservation_id = ?", reservationID).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, handlePaymentDBError(err)
	}
	return payments, nil
}

func (r *paymentRepository) FinalizePayment(ctx context.Context, id, userID int64, allowAny, succeeded bool, referenceID, failureReason string, now time.Time) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		if err := expireReservations(tx, now); err != nil {
			return handlePaymentDBError(err)
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrPaymentNotFound
			}
			return handlePaymentDBError(err)
		}
		if !allowAny && payment.UserID != userID {
			return exception.ErrNotAccessible
		}
		if payment.Status != enum.PaymentPending {
			return nil
		}

		var reservationRef entity.Reservation
		if err := tx.Select("id, showtime_id").First(&reservationRef, payment.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&showtime, reservationRef.ShowtimeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handlePaymentDBError(err)
		}
		var reservation entity.Reservation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, payment.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		if showtime.Status != enum.ShowtimeStatusPublished || !showtime.StartsAt.After(now) {
			return failPayment(tx, &payment, "SHOWTIME_NOT_BOOKABLE", now)
		}

		if reservation.Status == enum.ReservationExpired || !reservation.ExpiresAt.After(now) {
			return failPayment(tx, &payment, "RESERVATION_EXPIRED", now)
		}
		if reservation.Status != enum.ReservationPendingPayment || payment.Amount != reservation.TotalAmount {
			return failPayment(tx, &payment, "RESERVATION_NOT_PAYABLE", now)
		}
		if !succeeded {
			if strings.TrimSpace(failureReason) == "" {
				failureReason = "MOCK_PAYMENT_REJECTED"
			}
			return failPayment(tx, &payment, failureReason, now)
		}

		var itemCount int64
		if err := tx.Model(&entity.ReservationItem{}).Where("reservation_id = ?", reservation.ID).Count(&itemCount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		result := tx.Model(&entity.ShowtimeSeat{}).
			Where("reservation_id = ? AND status = ? AND held_until > ?", reservation.ID, enum.ShowtimeSeatHeld, now).
			Updates(map[string]any{"status": enum.ShowtimeSeatSold, "held_until": nil, "updated_at": now})
		if result.Error != nil {
			return handlePaymentDBError(result.Error)
		}
		if itemCount == 0 || result.RowsAffected != itemCount {
			return exception.ErrReservationNotPayable
		}
		if err := tx.Model(&entity.Reservation{}).Where("id = ? AND status = ?", reservation.ID, enum.ReservationPendingPayment).
			Updates(map[string]any{"status": enum.ReservationPaid, "updated_at": now}).Error; err != nil {
			return handlePaymentDBError(err)
		}
		if err := tx.Model(&payment).Updates(map[string]any{
			"status": enum.PaymentSucceeded, "reference_id": referenceID,
			"failure_reason": nil, "paid_at": now, "updated_at": now,
		}).Error; err != nil {
			return handlePaymentDBError(err)
		}
		payment.Status = enum.PaymentSucceeded
		payment.ReferenceID = &referenceID
		payment.FailureReason = nil
		payment.PaidAt = &now
		payment.UpdatedAt = now
		ticket := entity.Ticket{
			ReservationID: reservation.ID,
			PaymentID:     payment.ID,
			UserID:        payment.UserID,
			Code:          "ct_" + strings.ReplaceAll(uuid.NewString(), "-", ""),
			Status:        enum.TicketActive,
		}
		if err := tx.Create(&ticket).Error; err != nil {
			return handlePaymentDBError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) RefundPayment(ctx context.Context, id, userID int64, allowAny bool, referenceID string, now time.Time, cutoff time.Duration) (*entity.Payment, error) {
	if id <= 0 || userID <= 0 || strings.TrimSpace(referenceID) == "" || cutoff < 0 {
		return nil, exception.ErrInvalidInput
	}
	var payment entity.Payment
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)

		var paymentRef entity.Payment
		if err := tx.Select("id, reservation_id").First(&paymentRef, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrPaymentNotFound
			}
			return handlePaymentDBError(err)
		}
		var reservationRef entity.Reservation
		if err := tx.Select("id, showtime_id").First(&reservationRef, paymentRef.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&showtime, reservationRef.ShowtimeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handlePaymentDBError(err)
		}
		var reservation entity.Reservation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, paymentRef.ReservationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handlePaymentDBError(err)
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrPaymentNotFound
			}
			return handlePaymentDBError(err)
		}
		if !allowAny && payment.UserID != userID {
			return exception.ErrNotAccessible
		}
		if payment.Status == enum.PaymentRefunded {
			return nil
		}
		if payment.Status != enum.PaymentSucceeded || reservation.Status != enum.ReservationPaid {
			return exception.ErrRefundNotAllowed
		}
		if !now.Before(showtime.StartsAt.Add(-cutoff)) {
			return exception.ErrRefundNotAllowed
		}

		var ticket entity.Ticket
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("payment_id = ?", payment.ID).First(&ticket).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrTicketNotFound
			}
			return handlePaymentDBError(err)
		}
		if ticket.Status != enum.TicketActive {
			return exception.ErrRefundNotAllowed
		}

		var itemCount int64
		if err := tx.Model(&entity.ReservationItem{}).Where("reservation_id = ?", reservation.ID).Count(&itemCount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		seatResult := tx.Model(&entity.ShowtimeSeat{}).
			Where("reservation_id = ? AND status = ?", reservation.ID, enum.ShowtimeSeatSold).
			Updates(map[string]any{"status": enum.ShowtimeSeatAvailable, "reservation_id": nil, "held_until": nil, "updated_at": now})
		if seatResult.Error != nil {
			return handlePaymentDBError(seatResult.Error)
		}
		if itemCount == 0 || seatResult.RowsAffected != itemCount {
			return exception.ErrRefundNotAllowed
		}
		reservationResult := tx.Model(&entity.Reservation{}).
			Where("id = ? AND status = ?", reservation.ID, enum.ReservationPaid).
			Updates(map[string]any{"status": enum.ReservationCancelled, "updated_at": now})
		if reservationResult.Error != nil {
			return handlePaymentDBError(reservationResult.Error)
		}
		if reservationResult.RowsAffected != 1 {
			return exception.ErrRefundNotAllowed
		}
		ticketResult := tx.Model(&entity.Ticket{}).
			Where("id = ? AND status = ?", ticket.ID, enum.TicketActive).
			Updates(map[string]any{"status": enum.TicketCancelled, "updated_at": now})
		if ticketResult.Error != nil {
			return handlePaymentDBError(ticketResult.Error)
		}
		if ticketResult.RowsAffected != 1 {
			return exception.ErrRefundNotAllowed
		}
		paymentResult := tx.Model(&entity.Payment{}).
			Where("id = ? AND status = ?", payment.ID, enum.PaymentSucceeded).
			Updates(map[string]any{
				"status": enum.PaymentRefunded, "refund_reference_id": referenceID,
				"refunded_at": now, "updated_at": now,
			})
		if paymentResult.Error != nil {
			return handlePaymentDBError(paymentResult.Error)
		}
		if paymentResult.RowsAffected != 1 {
			return exception.ErrRefundNotAllowed
		}
		payment.Status = enum.PaymentRefunded
		payment.RefundReferenceID = &referenceID
		payment.RefundedAt = &now
		payment.UpdatedAt = now
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func failPayment(tx *gorm.DB, payment *entity.Payment, reason string, now time.Time) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "PAYMENT_FAILED"
	}
	runes := []rune(reason)
	if len(runes) > 255 {
		reason = string(runes[:255])
	}
	if err := tx.Model(payment).Updates(map[string]any{
		"status": enum.PaymentFailed, "failure_reason": reason, "updated_at": now,
	}).Error; err != nil {
		return handlePaymentDBError(err)
	}
	payment.Status = enum.PaymentFailed
	payment.FailureReason = &reason
	payment.UpdatedAt = now
	return nil
}

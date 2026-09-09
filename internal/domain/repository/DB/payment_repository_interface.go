package DB

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, bool, error)
	GetPayment(ctx context.Context, id int64) (*entity.Payment, error)
	ListReservationPayments(ctx context.Context, reservationID int64, limit, offset int) ([]entity.Payment, int64, error)
	FinalizePayment(ctx context.Context, id, userID int64, allowAny, succeeded bool, referenceID, failureReason string, now time.Time) (*entity.Payment, error)
	RefundPayment(ctx context.Context, id, userID int64, allowAny bool, referenceID string, now time.Time, cutoff time.Duration) (*entity.Payment, error)
}

package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type PaymentServiceInterface interface {
	Create(ctx context.Context, reservationID, userID int64, idempotencyKey string) (*dto.PaymentResponse, bool, error)
	Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.PaymentResponse, error)
	ListForReservation(ctx context.Context, reservationID, userID int64, allowAny bool, page, limit int) ([]*dto.PaymentResponse, int64, error)
	Verify(ctx context.Context, id, userID int64, allowAny bool, req *dto.VerifyMockPaymentRequest) (*dto.PaymentResponse, error)
	Refund(ctx context.Context, id, userID int64, allowAny bool) (*dto.PaymentResponse, error)
}

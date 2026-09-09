package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type TicketServiceInterface interface {
	Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.TicketResponse, error)
	ListMine(ctx context.Context, userID int64, page, limit int) ([]*dto.TicketResponse, int64, error)
	QRCode(ctx context.Context, id, userID int64, allowAny bool) ([]byte, error)
	PDF(ctx context.Context, id, userID int64, allowAny bool) (*dto.TicketPDFResponse, error)
	EnsurePDFForPayment(ctx context.Context, paymentID int64) error
	Validate(ctx context.Context, code string, validatorUserID int64, allowAny bool) (*dto.TicketResponse, error)
}

type TicketDocumentIssuer interface {
	EnsurePDFForPayment(ctx context.Context, paymentID int64) error
}

package DB

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type SalesReport struct {
	CinemaID           int64     `json:"cinema_id"`
	From               time.Time `json:"from"`
	To                 time.Time `json:"to"`
	GrossAmount        int64     `json:"gross_amount"`
	RefundedAmount     int64     `json:"refunded_amount"`
	NetAmount          int64     `json:"net_amount"`
	SuccessfulPayments int64     `json:"successful_payments"`
	RefundedPayments   int64     `json:"refunded_payments"`
	TicketsSold        int64     `json:"tickets_sold"`
}

type AuditFilter struct {
	ActorUserID *int64
	From        *time.Time
	To          *time.Time
	Limit       int
	Offset      int
}

type OperationsRepositoryInterface interface {
	SalesReport(ctx context.Context, cinemaID int64, from, to time.Time) (*SalesReport, error)
	CreateSettlement(ctx context.Context, cinemaID, createdByUserID int64, from, to time.Time) (*entity.CinemaSettlement, error)
	ListSettlements(ctx context.Context, cinemaID int64, limit, offset int) ([]entity.CinemaSettlement, int64, error)
	MarkSettlementPaid(ctx context.Context, settlementID, paidByUserID int64, now time.Time) (*entity.CinemaSettlement, error)
	RecordAudit(ctx context.Context, entry *entity.AuditLog) error
	ListAudit(ctx context.Context, filter AuditFilter) ([]entity.AuditLog, int64, error)
}

package service_contract

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type OperationsServiceInterface interface {
	SalesReport(ctx context.Context, cinemaID int64, from, to time.Time) (*DB.SalesReport, error)
	CreateSettlement(ctx context.Context, cinemaID, actorUserID int64, req *dto.CreateSettlementRequest) (*entity.CinemaSettlement, error)
	ListSettlements(ctx context.Context, cinemaID int64, page, limit int) ([]entity.CinemaSettlement, int64, error)
	MarkSettlementPaid(ctx context.Context, settlementID, actorUserID int64) (*entity.CinemaSettlement, error)
	ListAudit(ctx context.Context, filter dto.OperationsFilterRequest) ([]entity.AuditLog, int64, error)
}

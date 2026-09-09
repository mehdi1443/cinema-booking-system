package service

import (
	"context"
	"time"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type operationsService struct {
	repo DB.OperationsRepositoryInterface
}

func NewOperationsService(repo DB.OperationsRepositoryInterface) service_contract.OperationsServiceInterface {
	return &operationsService{repo: repo}
}

func validOperationsRange(from, to time.Time) bool {
	return !from.IsZero() && !to.IsZero() && to.After(from) && to.Sub(from) <= 366*24*time.Hour
}

func (s *operationsService) SalesReport(ctx context.Context, cinemaID int64, from, to time.Time) (*DB.SalesReport, error) {
	if cinemaID <= 0 || !validOperationsRange(from, to) {
		return nil, exception.ErrInvalidInput
	}
	return s.repo.SalesReport(ctx, cinemaID, from, to)
}

func (s *operationsService) CreateSettlement(ctx context.Context, cinemaID, actorUserID int64, req *dto.CreateSettlementRequest) (*entity.CinemaSettlement, error) {
	if cinemaID <= 0 || actorUserID <= 0 || req == nil || !validOperationsRange(req.From, req.To) || req.To.After(time.Now().Add(time.Minute)) {
		return nil, exception.ErrInvalidInput
	}
	return s.repo.CreateSettlement(ctx, cinemaID, actorUserID, req.From, req.To)
}

func (s *operationsService) ListSettlements(ctx context.Context, cinemaID int64, page, limit int) ([]entity.CinemaSettlement, int64, error) {
	if cinemaID <= 0 {
		return nil, 0, exception.ErrInvalidInput
	}
	_, limit, offset := normalizePagination(page, limit)
	return s.repo.ListSettlements(ctx, cinemaID, limit, offset)
}

func (s *operationsService) MarkSettlementPaid(ctx context.Context, settlementID, actorUserID int64) (*entity.CinemaSettlement, error) {
	if settlementID <= 0 || actorUserID <= 0 {
		return nil, exception.ErrInvalidInput
	}
	return s.repo.MarkSettlementPaid(ctx, settlementID, actorUserID, time.Now())
}

func (s *operationsService) ListAudit(ctx context.Context, filter dto.OperationsFilterRequest) ([]entity.AuditLog, int64, error) {
	if filter.From != nil && filter.To != nil && !validOperationsRange(*filter.From, *filter.To) {
		return nil, 0, exception.ErrInvalidInput
	}
	_, limit, offset := normalizePagination(filter.Page, filter.Limit)
	return s.repo.ListAudit(ctx, DB.AuditFilter{ActorUserID: filter.ActorUserID, From: filter.From, To: filter.To, Limit: limit, Offset: offset})
}

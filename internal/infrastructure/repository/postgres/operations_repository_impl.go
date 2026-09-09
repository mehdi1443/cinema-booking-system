package postgres_repo

import (
	"context"
	"errors"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type operationsRepository struct{ db database.Database }

func NewOperationsRepository(db database.Database) DB.OperationsRepositoryInterface {
	return &operationsRepository{db: db}
}

func (r *operationsRepository) ensureCinema(ctx context.Context, cinemaID int64) error {
	var count int64
	if err := r.db.GetDB().WithContext(ctx).Table("cinemas").Where("id = ? AND deleted_at IS NULL", cinemaID).Count(&count).Error; err != nil {
		return handlePaymentDBError(err)
	}
	if count == 0 {
		return exception.ErrCinemaNotFound
	}
	return nil
}

func (r *operationsRepository) SalesReport(ctx context.Context, cinemaID int64, from, to time.Time) (*DB.SalesReport, error) {
	if err := r.ensureCinema(ctx, cinemaID); err != nil {
		return nil, err
	}
	report := &DB.SalesReport{CinemaID: cinemaID, From: from, To: to}
	query := `SELECT
 COALESCE(SUM(CASE WHEN p.status IN ('SUCCEEDED','REFUNDED') THEN p.amount ELSE 0 END),0) gross_amount,
 COALESCE(SUM(CASE WHEN p.status = 'REFUNDED' THEN p.amount ELSE 0 END),0) refunded_amount,
 COALESCE(SUM(CASE WHEN p.status = 'SUCCEEDED' THEN p.amount ELSE 0 END),0) net_amount,
 COUNT(*) FILTER (WHERE p.status = 'SUCCEEDED') successful_payments,
 COUNT(*) FILTER (WHERE p.status = 'REFUNDED') refunded_payments,
 COUNT(DISTINCT t.id) FILTER (WHERE p.status = 'SUCCEEDED') tickets_sold
FROM payments p
JOIN reservations r ON r.id = p.reservation_id
JOIN showtimes s ON s.id = r.showtime_id
JOIN halls h ON h.id = s.hall_id
LEFT JOIN tickets t ON t.payment_id = p.id
WHERE h.cinema_id = ? AND p.paid_at >= ? AND p.paid_at < ?`
	if err := r.db.GetDB().WithContext(ctx).Raw(query, cinemaID, from, to).Scan(report).Error; err != nil {
		return nil, handlePaymentDBError(err)
	}
	return report, nil
}

func (r *operationsRepository) CreateSettlement(ctx context.Context, cinemaID, createdByUserID int64, from, to time.Time) (*entity.CinemaSettlement, error) {
	var settlement entity.CinemaSettlement
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var cinemaCount int64
		if err := tx.Table("cinemas").Where("id = ? AND deleted_at IS NULL", cinemaID).Count(&cinemaCount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		if cinemaCount == 0 {
			return exception.ErrCinemaNotFound
		}
		settlement = entity.CinemaSettlement{CinemaID: cinemaID, PeriodFrom: from, PeriodTo: to, Status: "PENDING", CreatedByUserID: createdByUserID}
		if err := tx.Create(&settlement).Error; err != nil {
			return handlePaymentDBError(err)
		}
		insert := `INSERT INTO cinema_settlement_items(settlement_id, payment_id, amount)
SELECT ?, p.id, p.amount
FROM payments p
JOIN reservations r ON r.id = p.reservation_id
JOIN showtimes s ON s.id = r.showtime_id
JOIN halls h ON h.id = s.hall_id
WHERE h.cinema_id = ? AND p.status = 'SUCCEEDED' AND p.paid_at >= ? AND p.paid_at < ? AND s.ends_at <= ?
  AND NOT EXISTS (SELECT 1 FROM cinema_settlement_items i WHERE i.payment_id = p.id)
FOR UPDATE OF p SKIP LOCKED`
		result := tx.Exec(insert, settlement.ID, cinemaID, from, to, to)
		if result.Error != nil {
			return handlePaymentDBError(result.Error)
		}
		if result.RowsAffected == 0 {
			return exception.ErrNoSettlementPayments
		}
		if err := tx.Table("cinema_settlement_items").Where("settlement_id = ?", settlement.ID).Select("COALESCE(SUM(amount),0)").Scan(&settlement.GrossAmount).Error; err != nil {
			return handlePaymentDBError(err)
		}
		return handlePaymentDBError(tx.Model(&settlement).Update("gross_amount", settlement.GrossAmount).Error)
	})
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

func (r *operationsRepository) ListSettlements(ctx context.Context, cinemaID int64, limit, offset int) ([]entity.CinemaSettlement, int64, error) {
	if err := r.ensureCinema(ctx, cinemaID); err != nil {
		return nil, 0, err
	}
	query := r.db.GetDB().WithContext(ctx).Model(&entity.CinemaSettlement{}).Where("cinema_id = ?", cinemaID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	var items []entity.CinemaSettlement
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	return items, total, nil
}

func (r *operationsRepository) MarkSettlementPaid(ctx context.Context, settlementID, paidByUserID int64, now time.Time) (*entity.CinemaSettlement, error) {
	var settlement entity.CinemaSettlement
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&settlement, settlementID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrSettlementNotFound
			}
			return handlePaymentDBError(err)
		}
		if settlement.Status == "PAID" {
			return nil
		}
		if err := tx.Model(&settlement).Updates(map[string]any{"status": "PAID", "paid_by_user_id": paidByUserID, "paid_at": now, "updated_at": now}).Error; err != nil {
			return handlePaymentDBError(err)
		}
		settlement.Status, settlement.PaidByUserID, settlement.PaidAt, settlement.UpdatedAt = "PAID", &paidByUserID, &now, now
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

func (r *operationsRepository) RecordAudit(ctx context.Context, entry *entity.AuditLog) error {
	if entry == nil {
		return exception.ErrInvalidInput
	}
	return handlePaymentDBError(r.db.GetDB().WithContext(ctx).Create(entry).Error)
}

func (r *operationsRepository) ListAudit(ctx context.Context, filter DB.AuditFilter) ([]entity.AuditLog, int64, error) {
	query := r.db.GetDB().WithContext(ctx).Model(&entity.AuditLog{})
	if filter.ActorUserID != nil {
		query = query.Where("actor_user_id = ?", *filter.ActorUserID)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at < ?", *filter.To)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	var items []entity.AuditLog
	if err := query.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, handlePaymentDBError(err)
	}
	return items, total, nil
}

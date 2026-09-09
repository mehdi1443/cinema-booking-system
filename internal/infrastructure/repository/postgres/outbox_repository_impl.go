package postgres_repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type outboxRepository struct{ db database.Database }

func NewOutboxRepository(db database.Database) DB.OutboxRepositoryInterface {
	return &outboxRepository{db: db}
}

func enqueueOutbox(tx *gorm.DB, eventType, aggregateType string, aggregateID int64, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return exception.Wrap(exception.ErrInternal, "failed to serialize outbox event", err)
	}
	event := entity.OutboxEvent{
		EventType: eventType, AggregateType: aggregateType, AggregateID: aggregateID,
		Payload: body, Status: "PENDING", AvailableAt: time.Now(),
	}
	return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "event_type"}, {Name: "aggregate_type"}, {Name: "aggregate_id"}}, DoNothing: true}).Create(&event).Error
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, now time.Time) ([]entity.OutboxEvent, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	var events []entity.OutboxEvent
	staleBefore := now.Add(-5 * time.Minute)
	query := `UPDATE outbox_events
SET status = 'PROCESSING', locked_at = ?, attempts = attempts + 1, updated_at = ?
WHERE id IN (
    SELECT id FROM outbox_events
    WHERE attempts < 10 AND (
        (status IN ('PENDING', 'FAILED') AND available_at <= ?)
        OR (status = 'PROCESSING' AND locked_at < ?)
    )
    ORDER BY available_at ASC, id ASC
    FOR UPDATE SKIP LOCKED
    LIMIT ?
)
RETURNING *`
	if err := r.db.GetDB().WithContext(ctx).Raw(query, now, now, now, staleBefore, limit).Scan(&events).Error; err != nil {
		return nil, handlePaymentDBError(err)
	}
	return events, nil
}

func (r *outboxRepository) MarkProcessed(ctx context.Context, id int64, now time.Time) error {
	result := r.db.GetDB().WithContext(ctx).Model(&entity.OutboxEvent{}).Where("id = ? AND status = ?", id, "PROCESSING").Updates(map[string]any{
		"status": "PROCESSED", "processed_at": now, "locked_at": nil, "last_error": nil, "updated_at": now,
	})
	if result.Error != nil {
		return handlePaymentDBError(result.Error)
	}
	return nil
}

func (r *outboxRepository) MarkFailed(ctx context.Context, id int64, reason string, retryAt time.Time) error {
	if len(reason) > 2000 {
		reason = reason[:2000]
	}
	result := r.db.GetDB().WithContext(ctx).Model(&entity.OutboxEvent{}).Where("id = ? AND status = ?", id, "PROCESSING").Updates(map[string]any{
		"status": "FAILED", "available_at": retryAt, "locked_at": nil, "last_error": reason, "updated_at": time.Now(),
	})
	if result.Error != nil {
		return handlePaymentDBError(result.Error)
	}
	return nil
}

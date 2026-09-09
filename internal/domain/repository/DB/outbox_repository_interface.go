package DB

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type OutboxRepositoryInterface interface {
	Claim(ctx context.Context, limit int, now time.Time) ([]entity.OutboxEvent, error)
	MarkProcessed(ctx context.Context, id int64, now time.Time) error
	MarkFailed(ctx context.Context, id int64, reason string, retryAt time.Time) error
}

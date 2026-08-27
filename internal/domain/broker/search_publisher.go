package broker

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type SearchEventPublisher interface {
	PublishSyncEvent(ctx context.Context, doc entity.SearchDocument) error
	PublishDeleteEvent(ctx context.Context, documentID string) error
}

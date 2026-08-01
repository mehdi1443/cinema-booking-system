package redis

import (
	"context"
	"time"
)

type RefreshTokenInterFace interface {
	Save(ctx context.Context, userID string, token string, ttl time.Duration) error
	Get(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, userID string) error
}

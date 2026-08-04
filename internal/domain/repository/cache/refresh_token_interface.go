package cache

import (
	"context"
	"time"
)

type RefreshTokenInterface interface {
	SaveRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}

package cache

import (
	"context"
	"time"
)

type TokenBlacklistInterface interface {
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (*bool, error)
}

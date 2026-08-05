package redis_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
)

type tokenBlacklistRepository struct {
	client database.Cache
}

func NewTokenBlacklistRedis(client database.Cache) cache.TokenBlacklistInterface {
	return &tokenBlacklistRepository{
		client: client,
	}
}

func (r *tokenBlacklistRepository) buildKey(token string) string {
	return fmt.Sprintf("auth:blacklist:%s", token)
}

func (r *tokenBlacklistRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	key := r.buildKey(token)
	err := r.client.GetRDB().Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		return fmt.Errorf("%w: failed to blacklist token: %v", exception.ErrCache, err)
	}
	return nil
}

func (r *tokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (*bool, error) {
	key := r.buildKey(token)
	exists, err := r.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to check token blacklist status: %v", exception.ErrCache, err)
	}
	success := exists > 0
	return &success, nil
}

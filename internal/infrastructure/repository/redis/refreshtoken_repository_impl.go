package redis_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository struct {
	client database.Cache
}

func NewRefreshTokenRepository(client database.Cache) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) getkeysignuture(userID string) string {
	//  token:refresh: REDIS_STRING
	return r.client.GetRedisConstant().RefreshTokenSignature + userID
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	key := r.getkeysignuture(userID)
	if err := r.client.GetRDB().Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("redis save refresh token error: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := r.getkeysignuture(userID)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", exception.ErrRefreshTokenNotFound
		}
		return "", fmt.Errorf("redis get refresh token error: %w", err)
	}
	return val, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := r.getkeysignuture(userID)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete refresh token error: %w", err)
	}
	return nil
}

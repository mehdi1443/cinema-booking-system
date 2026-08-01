package refresh_token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	redis_contract "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/redis"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository struct {
	client database.Cache
}

func NewRefreshTokenRepository(client database.Cache) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) getKey(userID string) string {
	return r.client.GetRedisConstant().RefreshTokenSignature + userID
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	key := r.getKey(userID)
	if err := r.client.GetRDB().Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("redis save refresh token error: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := r.getKey(userID)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", redis_contract.ErrRefreshTokenNotFound
		}
		return "", fmt.Errorf("redis get refresh token error: %w", err)
	}
	return val, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := r.getKey(userID)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete refresh token error: %w", err)
	}
	return nil
}

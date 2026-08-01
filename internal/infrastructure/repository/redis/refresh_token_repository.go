package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository struct {
	client *redis.Client
}

func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.client.Set(ctx, "refresh:"+userID, token, ttl).Err()
}

func (r *RefreshTokenRepository) Get(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, "refresh:"+userID).Result()
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, userID string) error {
	return r.client.Del(ctx, "refresh:"+userID).Err()
}

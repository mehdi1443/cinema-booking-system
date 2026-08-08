package redis_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type captchaRepository struct {
	client database.Cache
}

func NewCaptchaRepositoryRedis(client database.Cache) cache.CaptchaInterface {
	return &captchaRepository{
		client: client,
	}
}

func (r *captchaRepository) buildKey(captchaId string) string {
	return fmt.Sprintf("captcha:%s", captchaId)
}

func (r *captchaRepository) Save(ctx context.Context, captchaId string, code string, ttl time.Duration) error {
	key := r.buildKey(captchaId)
	if err := r.client.GetRDB().Set(ctx, key, code, ttl); err != nil {
		return fmt.Errorf("failed to save captcha [%s]: %w", captchaId, err)
	}
	return nil
}

func (r *captchaRepository) Get(ctx context.Context, captchaId string) (string, error) {
	key := r.buildKey(captchaId)
	code, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", exception.ErrCaptchaNotFound
		}
		return "", fmt.Errorf("failed to retrieve captcha [%s]: %w", captchaId, err)
	}
	return code, nil
}

func (r *captchaRepository) Delete(ctx context.Context, captchaId string) error {
	key := r.buildKey(captchaId)
	if err := r.client.GetRDB().Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete captcha [%s]: %w", captchaId, err)
	}
	return nil
}

func (r *captchaRepository) Exist(ctx context.Context, captchaId string) (*bool, error) {
	key := r.buildKey(captchaId)

	exists, err := r.client.GetRDB().Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check existence for captcha [%s]: %w", captchaId, err)
	}
	success := exists > 0

	return &success, nil
}

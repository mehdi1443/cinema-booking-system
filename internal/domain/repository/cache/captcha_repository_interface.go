package cache

import (
	"context"
	"time"
)

type CaptchaInterface interface {
	Save(ctx context.Context, captchaId string, code string, ttl time.Duration) error
	Get(ctx context.Context, captchaId string) (string, error)
	Delete(ctx context.Context, captchaId string) error
	Exist(ctx context.Context, captchaId string) (*bool, error)
}

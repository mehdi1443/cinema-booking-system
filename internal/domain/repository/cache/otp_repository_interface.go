package cache

import (
	"context"
	"time"
)

type OTPInterface interface {
	SaveOTP(ctx context.Context, tagId string, code string, ttl time.Duration) error
	GetOTP(ctx context.Context, tagId string) (string, error)
	ConsumeOTP(ctx context.Context, tagId string, code string) error
	DeleteOTP(ctx context.Context, tagId string) error
	SetVerified(ctx context.Context, tagId string, ttl time.Duration) error
	IsVerified(ctx context.Context, tagId string) (bool, error)
	DeleteVerified(ctx context.Context, tagId string) error
}

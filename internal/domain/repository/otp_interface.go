package repository_contract

import (
	"context"
	"errors"
	"time"
)

var ErrOtpKeyNotFound = errors.New("otp key not found in redis")

type OTPRepoInterface interface {
	SaveOTP(ctx context.Context, phone string, code string, ttl time.Duration) error
	GetOTP(ctx context.Context, phone string) (string, error)
	DeleteOTP(ctx context.Context, phone string) error
	SetVerified(ctx context.Context, phone string, ttl time.Duration) error
	IsVerified(ctx context.Context, phone string) (bool, error)
	DeleteVerified(ctx context.Context, phone string) error
}

package otp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	redis_contract "github.com/mehdi1443/cinema-booking-system/internal/infrastructure/repository/redis"

	"github.com/redis/go-redis/v9"
)

type OTPRepository struct {
	client database.Cache
}

func NewOTPRepository(client database.Cache) *OTPRepository {
	return &OTPRepository{client: client}
}

func (r *OTPRepository) getKey(prefix, phone string) string {
	return prefix + phone
}

func (r *OTPRepository) SaveOTP(ctx context.Context, phone string, code string, ttl time.Duration) error {
	key := r.getKey(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	return r.client.GetRDB().Set(ctx, key, code, ttl).Err()
}

func (r *OTPRepository) GetOTP(ctx context.Context, phone string) (string, error) {
	key := r.getKey(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", redis_contract.ErrOtpKeyNotFound
		}
		return "", fmt.Errorf("redis get otp error: %w", err)
	}
	return val, nil
}

func (r *OTPRepository) DeleteOTP(ctx context.Context, phone string) error {
	key := r.getKey(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	return r.client.GetRDB().Del(ctx, key).Err()
}

func (r *OTPRepository) SetVerified(ctx context.Context, phone string, ttl time.Duration) error {
	key := r.getKey(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
	return r.client.GetRDB().Set(ctx, key, "true", ttl).Err()
}

func (r *OTPRepository) IsVerified(ctx context.Context, phone string) (bool, error) {
	key := r.getKey(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("redis get verified status error: %w", err)
	}
	return val == "true", nil
}

func (r *OTPRepository) DeleteVerified(ctx context.Context, phone string) error {
	key := r.getKey(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
	return r.client.GetRDB().Del(ctx, key).Err()
}

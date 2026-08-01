package otp

import (
	"context"
	"errors"
	"fmt"
	"time"

	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type OTPRepository struct {
	client database.Cache
}

func NewOTPRepository(client database.Cache) *OTPRepository {
	return &OTPRepository{client: client}
}

func (r *OTPRepository) getkeysignuture(prefix, phone string) string {
	// otp:phone : REDIS_STRING
	// otp:verified : REDIS_STRING(false|true)
	return prefix + phone
}

func (r *OTPRepository) SaveOTP(ctx context.Context, phone string, code string, ttl time.Duration) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	if err := r.client.GetRDB().Set(ctx, key, code, ttl).Err(); err != nil {
		return fmt.Errorf("redis save otp error: %w", err)
	}
	return nil
}

func (r *OTPRepository) GetOTP(ctx context.Context, phone string) (string, error) {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", repository_contract.ErrOtpKeyNotFound
		}
		return "", fmt.Errorf("redis get otp error: %w", err)
	}
	return val, nil
}

func (r *OTPRepository) DeleteOTP(ctx context.Context, phone string) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, phone)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete otp error: %w", err)
	}
	return nil
}

func (r *OTPRepository) SetVerified(ctx context.Context, phone string, ttl time.Duration) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
	if err := r.client.GetRDB().Set(ctx, key, "true", ttl).Err(); err != nil {
		return fmt.Errorf("redis set verified status error: %w", err)
	}
	return nil
}

func (r *OTPRepository) IsVerified(ctx context.Context, phone string) (bool, error) {
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
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
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, phone)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete verified status error: %w", err)
	}
	return nil
}

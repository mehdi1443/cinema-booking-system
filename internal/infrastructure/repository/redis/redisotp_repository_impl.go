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

type OTPRepository struct {
	client database.Cache
}

func NewOTPRepository(client database.Cache) *OTPRepository {
	return &OTPRepository{client: client}
}

func (r *OTPRepository) getkeysignuture(prefix, tagId string) string {
	// otp:phone_email:tag|email : REDIS_STRING
	// otp:phone_email:verified : REDIS_STRING(false|true)
	return prefix + tagId
}

func (r *OTPRepository) SaveOTP(ctx context.Context, tagId string, code string, ttl time.Duration) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, tagId)
	if err := r.client.GetRDB().Set(ctx, key, code, ttl).Err(); err != nil {
		return fmt.Errorf("redis save otp error: %w", err)
	}
	return nil
}

func (r *OTPRepository) GetOTP(ctx context.Context, tagId string) (string, error) {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, tagId)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", exception.ErrOtpKeyNotFound
		}
		return "", fmt.Errorf("redis get otp error: %w", err)
	}
	return val, nil
}

func (r *OTPRepository) DeleteOTP(ctx context.Context, tagId string) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().PhoneOtpSignature, tagId)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete otp error: %w", err)
	}
	return nil
}

func (r *OTPRepository) SetVerified(ctx context.Context, tagId string, ttl time.Duration) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, tagId)
	if err := r.client.GetRDB().Set(ctx, key, "true", ttl).Err(); err != nil {
		return fmt.Errorf("redis set verified status error: %w", err)
	}
	return nil
}

func (r *OTPRepository) IsVerified(ctx context.Context, tagId string) (bool, error) {
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, tagId)
	val, err := r.client.GetRDB().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("redis get verified status error: %w", err)
	}
	return val == "true", nil
}

func (r *OTPRepository) DeleteVerified(ctx context.Context, tagId string) error {
	key := r.getkeysignuture(r.client.GetRedisConstant().VerifiedOtpSignature, tagId)
	if err := r.client.GetRDB().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete verified status error: %w", err)
	}
	return nil
}

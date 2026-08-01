package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPRepository struct {
	client *redis.Client
}

func NewOTPRepository(client *redis.Client) *OTPRepository {
	return &OTPRepository{client: client}
}

func (r *OTPRepository) Save(phone string, code string, ttl time.Duration) error {
	return r.client.Set(context.Background(), "otp:"+phone, code, ttl).Err()
}

func (r *OTPRepository) Get(phone string) (string, error) {
	return r.client.Get(context.Background(), "otp:"+phone).Result()
}

func (r *OTPRepository) Delete(phone string) error {
	return r.client.Del(context.Background(), "otp:"+phone).Err()
}

func (r *OTPRepository) SetVerified(phone string, ttl time.Duration) error {
	return r.client.Set(context.Background(), "verified:"+phone, "true", ttl).Err()
}

func (r *OTPRepository) IsVerified(phone string) (bool, error) {
	val, err := r.client.Get(context.Background(), "verified:"+phone).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "true", err
}

func (r *OTPRepository) DeleteVerified(phone string) error {
	return r.client.Del(context.Background(), "verified:"+phone).Err()
}

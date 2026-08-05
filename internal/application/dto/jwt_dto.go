package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type GenerateTokenRequest struct {
	UserID int64 `json:"user_id"`
}

type TokenInContext struct {
	UserID int64         `json:"user_id"`
	Roles  []entity.Role `json:"roles"`
}

type JWTTokenInContext struct {
	TheToken         string    `json:"the_token"`
	TheTokenExpireAt time.Time `json:"the_token_expire_at"`
}

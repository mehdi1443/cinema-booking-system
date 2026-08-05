package dto

import (
	"time"
)

type GenerateTokenRequest struct {
	UserID int64 `json:"user_id"`
}

type TokenInContext struct {
	UserID          int64    `json:"user_id"`
	PermissionNames []string `json:"permission_names"`
}

type JWTTokenInContext struct {
	TheToken         string    `json:"the_token"`
	TheTokenExpireAt time.Time `json:"the_token_expire_at"`
}

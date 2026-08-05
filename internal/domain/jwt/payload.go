package domainjwt

import (
	"time"
)

type TokenType string

const (
	AccessTokenTag  TokenType = "access"
	RefreshTokenTag TokenType = "refresh"
)

const (
	CtxUserIDKey           = "user_id"
	CtxRolesKey            = "roles"
	CtxPayloadKey          = "token_payload"
	CtxTheTokenKey         = "the_token"
	CtxTheTokenExpireAtKey = "the_token_expire_at"
)

type TokenPayload struct {
	UserID    int64     `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	IssuedAt  time.Time `json:"issued_at"`  // Token issue time (iat)
	ExpiresAt time.Time `json:"expired_at"` // Token expiration time (exp)
}

func (p *TokenPayload) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

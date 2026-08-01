package domainjwt

import (
	"time"
)

type TokenType string

const (
	AccessTokenTag  TokenType = "access"
	RefreshTokenTag TokenType = "refresh"
)

type TokenPayload struct {
	UserID    int64     `json:"user_id"`
	Roles     []string  `json:"roles"`
	TokenType TokenType `json:"token_type"`
	IssuedAt  time.Time `json:"issued_at"`  // Token issue time (iat)
	ExpiresAt time.Time `json:"expired_at"` // Token expiration time (exp)
}

func (p *TokenPayload) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

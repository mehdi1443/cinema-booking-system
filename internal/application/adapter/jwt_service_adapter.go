package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

const (
	UserIDKeyMapClaims    string = "user_id"
	TokenTypeKeyMapClaims string = "token_type"
	IssuedAtKeyMapClaims  string = "iat"
	ExpireAtKeyMapClaims  string = "exp"
)

type JWTServiceInterface interface {
	GenerateTokens(ctx context.Context, req dto.GenerateTokenRequest) (access string, refresh string, err error)
	GenerateAccessToken(ctx context.Context, req dto.GenerateTokenRequest) (string, error)
	GenerateRefreshToken(ctx context.Context, req dto.GenerateTokenRequest) (string, error)
	ValidateToken(ctx context.Context, signedToken string, tokenType string) (*domainjwt.TokenPayload, error)
}

package service_contract

import (
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
	GenerateAccessToken(req dto.GenerateTokenRequest) (string, error)
	GenerateRefreshToken(req dto.GenerateTokenRequest) (string, error)
	ValidateToken(signedToken string) (*domainjwt.TokenPayload, error)
}

package service_contract

import (
	"errors"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid or malformed token")
)

const (
	UserIDKeyMapClaims    string = "user_id"
	TokenTypeKeyMapClaims string = "token_type"
	IssuedAtKeyMapClaims  string = "iat"
	ExpireAtKeyMapClaims  string = "exp"
)

type (
	JWTServiceInterface interface {
		GenerateAccessToken(req dto.GenerateTokenRequest) (string, error)
		GenerateRefreshToken(req dto.GenerateTokenRequest) (string, error)
		ValidateToken(signedToken string) (*domainjwt.TokenPayload, error)
	}

	UserServiceInterface interface {
	}
)

package service_contract

import (
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
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

package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type AuthenticationServiceInterface interface {
	SendOTP(ctx context.Context, phone string, email string) (*bool, error)
	VerifyOTP(ctx context.Context, phone string, code string) (*bool, error)

	Login(ctx context.Context, phone string, password string) (*dto.StartAuthResponse, error)
	Register(ctx context.Context, phone string, email string, password string) (*dto.StartAuthResponse, error)

	VerifyPhone(ctx context.Context, userId int64, phone string, code string) (*bool, error)
	VerifyEmail(ctx context.Context, tokenData *dto.TokenInContext, email string, code string) (*bool, error)

	RefreshToken(ctx context.Context, refreshToken *dto.RefreshTokenRequest) (*dto.StartAuthResponse, error)
	Logout(ctx context.Context, userID int64, jwtContext *dto.JWTTokenInContext, refreshToken *dto.LogoutRequest) error
}

package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type AuthenticationServiceInterface interface {
	SendOTP(ctx context.Context, phone string) (*dto.SendOTPResponse, error)
	VerifyOTP(ctx context.Context, phone string, code string) (*dto.StartAuthResponse, error)

	RefreshToken(ctx context.Context, refreshToken *dto.RefreshTokenRequest) (*dto.StartAuthResponse, error)
	Logout(ctx context.Context, userID int64, jwtContext *dto.JWTTokenInContext, refreshToken *dto.LogoutRequest) error
}

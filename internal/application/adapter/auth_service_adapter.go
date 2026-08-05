package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type AuthenticationServiceInterface interface {
	SendOTP(ctx context.Context, phone string) (*bool, error)
	VerifyOTP(ctx context.Context, phone string, code string) (*bool, error)

	Login(ctx context.Context, phone string, password string) (access string, refresh string, err error)
	Register(ctx context.Context, phone string, email string, password string) (access string, refresh string, err error)

	VerifyPhone(ctx context.Context, tokenData *dto.TokenInContext, phone string, code string) (*bool, error)
	VerifyEmail(ctx context.Context, tokenData *dto.TokenInContext, email string, code string) (*bool, error)
}

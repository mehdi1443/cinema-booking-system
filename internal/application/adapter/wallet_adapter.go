package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type WalletServiceInterface interface {
	GetBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64) (int64, error)
	ChangeBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64, amount int64) (int64, error)
}

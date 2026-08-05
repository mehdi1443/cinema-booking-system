package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type walletService struct {
	userRepo   repository_contract.UserInterface
	walletRepo repository_contract.WalletInterface
}

func NewWalletService(userRepo repository_contract.UserInterface, walletRepo repository_contract.WalletInterface) service_contract.WalletServiceInterface {
	return &walletService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

// GetBalance retrieves the user's current wallet balance via the repository.
func (s *walletService) GetBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64) (int64, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

// ChangeBalance validates the amount, checks available balance, and delegates execution to the repository.
func (s *walletService) ChangeBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64, amount int64) (int64, error) {
	if amount == 0 {
		return 0, exception.ErrInvalidAmount
	}

	currentBalance, err := s.GetBalance(ctx, tokenData, userID)
	if err != nil {
		return 0, err
	}

	if currentBalance+amount < 0 {
		return 0, exception.ErrInsufficientBalance
	} else if currentBalance+amount > entity.MaxWalletBalance {
		return 0, exception.ErrInsufficientBalance
	}

	return s.walletRepo.AddAmountToBalance(ctx, userID, amount)

}

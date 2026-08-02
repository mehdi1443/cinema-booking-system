package wallet

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
)

type walletService struct {
	userRepo repository_contract.UserInterFace
	db       database.Database
}

func NewTokenService() service_contract.WalletServiceInterface {
	return &walletService{}
}

func NewWalletService(userRepo repository_contract.UserInterFace, db database.Database) service_contract.WalletServiceInterface {
	return &walletService{
		userRepo: userRepo,
		db:       db,
	}
}

// GetBalance retrieves the user's current wallet balance via the repository.
func (s *walletService) GetBalance(ctx context.Context, userID int64) (int64, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

// IncreaseBallance validates the input amount and delegates execution to the repository.
func (s *walletService) IncreaseBallance(ctx context.Context, userID int64, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, repository_contract.ErrInvalidAmount
	}

	return s.userRepo.IncreaseBallance(ctx, userID, amount)
}

// DecreaseBalance validates the amount, checks available balance, and delegates execution to the repository.
func (s *walletService) DecreaseBalance(ctx context.Context, userID int64, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, repository_contract.ErrInvalidAmount
	}

	currentBalance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return 0, err
	}

	if currentBalance < amount {
		return 0, repository_contract.ErrInsufficientBalance
	}

	return s.userRepo.DecreaseBalance(ctx, userID, amount)
}

package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type walletService struct {
	userRepo   repository_contract.UserInterface
	walletRepo repository_contract.WalletInterface
	logger     domainlogger.Logger
}

func NewWalletService(
	userRepo repository_contract.UserInterface,
	walletRepo repository_contract.WalletInterface,
	log domainlogger.Logger,
) service_contract.WalletServiceInterface {
	return &walletService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		logger:     log,
	}
}

func (s *walletService) GetBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64) (int64, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

func (s *walletService) ChangeBalance(ctx context.Context, tokenData *dto.TokenInContext, userID int64, amount int64) (int64, error) {
	if amount == 0 {
		return 0, exception.ErrInvalidAmount
	}

	currentBalance, err := s.GetBalance(ctx, tokenData, userID)
	if err != nil {
		return 0, err
	}

	// Check the bounds before adding so a malicious int64 amount cannot wrap
	// around and bypass the wallet limit.
	if currentBalance < 0 || currentBalance > entity.MaxWalletBalance ||
		(amount < 0 && amount < -currentBalance) {
		s.logger.Warn(ctx, "invalid balance adjustment request", domainlogger.Error(exception.ErrInsufficientBalance))
		return 0, exception.ErrInsufficientBalance
	}
	if amount > 0 && amount > entity.MaxWalletBalance-currentBalance {
		return 0, exception.ErrReachManBalance
	}

	return s.walletRepo.AddAmountToBalance(ctx, userID, amount)
}

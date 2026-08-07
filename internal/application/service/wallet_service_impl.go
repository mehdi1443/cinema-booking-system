package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

type walletService struct {
	userRepo   repository_contract.UserInterface
	walletRepo repository_contract.WalletInterface
	logger     logger.Logger
}

func NewWalletService(
	userRepo repository_contract.UserInterface,
	walletRepo repository_contract.WalletInterface,
	log logger.Logger,
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

	if currentBalance+amount < 0 || currentBalance+amount > entity.MaxWalletBalance {
		s.logger.Warn(logger.Validation, logger.Api, "invalid balance adjustment request", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: exception.ErrInsufficientBalance.Error(),
		})
		return 0, exception.ErrInsufficientBalance
	}

	return s.walletRepo.AddAmountToBalance(ctx, userID, amount)
}

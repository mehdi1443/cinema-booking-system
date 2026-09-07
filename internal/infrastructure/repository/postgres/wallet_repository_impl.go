package postgres_repo

import (
	"context"
	"errors"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type walletRepositoryPostgres struct {
	db database.Database
}

func NewWalletRepositoryPostgres(db database.Database) repository_contract.WalletInterface {
	return &walletRepositoryPostgres{
		db: db,
	}
}

func (r *walletRepositoryPostgres) AddAmountToBalance(
	ctx context.Context,
	userID int64,
	amount int64,
) (int64, error) {

	var finalBalance int64

	err := r.db.
		GetDB().
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			var user entity.User

			err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND deleted_at IS NULL", userID).
				First(&user).
				Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return exception.ErrUserNotFound
				}
				return handleDBError(err)
			}

			if user.WalletBalance+amount < 0 {
				return exception.ErrInsufficientBalance
			} else if user.WalletBalance+amount > entity.MaxWalletBalance {
				return exception.ErrReachManBalance
			}

			finalBalance = user.WalletBalance + amount

			err = tx.
				Model(&user).
				Update("wallet_balance", finalBalance).
				Error

			if err != nil {
				return handleDBError(err)
			}

			return nil
		})

	if err != nil {
		return 0, err
	}

	return finalBalance, nil
}

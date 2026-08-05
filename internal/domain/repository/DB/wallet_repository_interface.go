package DB

import "context"

type WalletInterface interface {
	AddAmountToBalance(ctx context.Context, userID int64, amount int64) (int64, error)
}

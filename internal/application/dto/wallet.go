package dto

type ChangeBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

type WalletBalanceResponse struct {
	UserID  int64 `json:"user_id"`
	Balance int64 `json:"balance"`
}

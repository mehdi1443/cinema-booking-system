package dto

type GenerateTokenRequest struct {
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

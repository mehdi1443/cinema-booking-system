package domainjwt

type TokenPayload struct {
	UserID    int64    `json:"user_id"`
	Roles     []string `json:"roles"`
	TokenType string   `json:"token_type"`
	ExpiresAt int64    `json:"expire_at"`
	CreatedAt int64    `json:"created_at"`
}

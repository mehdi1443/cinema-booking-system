package dto

type SendEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Text  string `json:"text"`
}

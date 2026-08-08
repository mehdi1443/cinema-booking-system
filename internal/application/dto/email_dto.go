package dto

type SendEmailRequest struct {
	To          []string `json:"to" binding:"required,email"`
	Subject     string   `json:"subject" binding:"required"`
	TextBody    string   `json:"text_body,omitempty"`
	HTMLBody    string   `json:"html_body,omitempty"`
	EmbedFiles  []string `json:"embed_files,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type SendOTPEmailRequest struct {
	To               []string `json:"to" binding:"required,email"`
	OTPCode          string   `json:"otp-code" binding:"required,otp-code"`
	OTPExpireMinutes int      `json:"otp-expr-minutes" binding:"required"`
}

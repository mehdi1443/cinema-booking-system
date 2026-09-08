package dto

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required,iranian-phone"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required,iranian-phone"`
	Code  string `json:"code" binding:"required,otp-code"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,iranian-phone"`
	Password string `json:"password" binding:"required,user-password"`
}

type StartAuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	UserID       int64    `json:"user_id"`
	Phone        string   `json:"phone"`
	Email        *string  `json:"email,omitempty"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	IsNewUser    bool     `json:"is_new_user"`
	Roles        []string `json:"roles"`
}

type SendOTPResponse struct {
	Phone            string `json:"phone"`
	OTPCode          string `json:"otp_code"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,iranian-phone"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,user-password"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required,iranian-phone"`
	Code  string `json:"code" binding:"required,otp-code"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,otp-code"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

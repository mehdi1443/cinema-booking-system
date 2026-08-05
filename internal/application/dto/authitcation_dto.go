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

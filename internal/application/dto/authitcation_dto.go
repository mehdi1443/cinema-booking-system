package dto

type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required"` // TODO : Vallation
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"` 
	Code  string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"` // TODO : Vallation
	Email    string `json:"email" binding:"required"` // TODO : Vallation
	Password string `json:"password" binding:"required,min=8"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

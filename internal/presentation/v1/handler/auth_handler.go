package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type AuthHandler struct {
	authService         service_contract.AuthenticationServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
}

func NewAuthHandler(authService service_contract.AuthenticationServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// TODO : Check Switch Cases Completeness
// SendOTP handles POST /auth/otp/send
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.SendOTP(c.Request.Context(), req.Phone)
	if err != nil {
		log.Printf("[SERVICE ERROR] SendOTP failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidPhone):
			response.Error(c, http.StatusBadRequest, "Invalid phone number", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to send OTP", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "OTP sent successfully", gin.H{
		"success": success,
	})
}

// VerifyOTP handles POST /auth/otp/verify
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		log.Printf("[SERVICE ERROR] VerifyOTP failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			response.Error(c, http.StatusBadRequest, "OTP code has expired", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to verify OTP", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", gin.H{
		"success": success,
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	access, ref, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		log.Printf("[SERVICE ERROR] Login failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "Invalid phone number or password", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Login failed", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Login successful", dto.AuthTokenResponse{
		AccessToken:  access,
		RefreshToken: ref,
	})
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	access, ref, err := h.authService.Register(c.Request.Context(), req.Phone, req.Email, req.Password)
	if err != nil {
		log.Printf("[SERVICE ERROR] Register failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, "User with this phone number already exists", err)
		case errors.Is(err, exception.ErrInvalidPhone):
			response.Error(c, http.StatusBadRequest, "Invalid phone number", err)
		case errors.Is(err, exception.ErrInvalidPassword):
			response.Error(c, http.StatusBadRequest, "Invalid password", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Registration failed", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Registration successful", dto.AuthTokenResponse{
		AccessToken:  access,
		RefreshToken: ref,
	})
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {

}

func (h *AuthHandler) Logout(ctx *gin.Context) {

}

// VerifyPhone handles POST /auth/verify-phone
func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyPhone(c.Request.Context(), &tokenData, req.Phone, req.Code)
	if err != nil {
		log.Printf("[SERVICE ERROR] VerifyPhone failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			response.Error(c, http.StatusBadRequest, "OTP code has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to verify phone number", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Phone number verified successfully", gin.H{
		"success": success,
	})
}

// VerifyEmail handles POST /auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyEmail(c.Request.Context(), &tokenData, req.Email, req.Code)
	if err != nil {
		log.Printf("[SERVICE ERROR] VerifyEmail failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			response.Error(c, http.StatusBadRequest, "OTP code has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to verify email address", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Email address verified successfully", gin.H{
		"success": success,
	})
}

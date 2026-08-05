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

func NewAuthHandler(authService service_contract.AuthenticationServiceInterface, contexttokenService service_contract.ContextTokenServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		contexttokenService: contexttokenService,
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
	log.Println(req)
	if (req.Email == req.Phone && req.Email == "") || (req.Email != "" && req.Phone != "") {
		response.Error(c, http.StatusBadRequest, "Just fill extract one field(email or phone)")
		return
	}

	success, err := h.authService.SendOTP(c.Request.Context(), req.Phone, req.Email)
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
			response.Error(c, http.StatusBadRequest, "OTP code has expired or phone not found", err)
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

	res, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		log.Printf("[SERVICE ERROR] Login failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrNotVerifiedPhoneNumber):
			response.Error(c, http.StatusUnauthorized, "Your phone number has not verified yet", err)
		case errors.Is(err, exception.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, "Invalid phone number or password", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Login failed", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req.Phone, req.Email, req.Password)
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

	response.Success(c, http.StatusCreated, "Registration successful", res)
}

// Refresh handles POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] Refresh token failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidToken):
			response.Error(c, http.StatusUnauthorized, "Invalid refresh token", err)
		case errors.Is(err, exception.ErrTokenExpired):
			response.Error(c, http.StatusUnauthorized, "Refresh token has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to refresh token", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", tokens)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	jwtData, err := h.contexttokenService.GetJWTFromGinContext(c)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to extract jwt context: %v", err)
		response.Error(c, http.StatusUnauthorized, exception.ErrGetJWTDataFromContext.Error())
		return
	}

	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), &tokenData, &jwtData, &req); err != nil {
		log.Printf("[SERVICE ERROR] Logout failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidToken):
			response.Error(c, http.StatusBadRequest, "Invalid token", err)
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to logout", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
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

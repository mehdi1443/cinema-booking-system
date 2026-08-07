package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type AuthHandler struct {
	authService         service_contract.AuthenticationServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewAuthHandler(
	authService service_contract.AuthenticationServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// TODO : Check Switch Cases Completeness
// SendOTP handles POST /auth/otp/send
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for SendOTP", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
			logger.Path:         c.Request.URL.Path,
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if (req.Email == req.Phone && req.Email == "") || (req.Email != "" && req.Phone != "") {
		h.log.Warn(logger.Validation, logger.Api, "SendOTP request contained both or neither email/phone", map[logger.ExtraKey]interface{}{
			logger.Path: c.Request.URL.Path,
		})
		response.Error(c, http.StatusBadRequest, "Just fill extract one field(email or phone)")
		return
	}

	success, err := h.authService.SendOTP(c.Request.Context(), req.Phone, req.Email)
	if err != nil {
		if errors.Is(err, exception.ErrInvalidPhone) {
			h.log.Warn(logger.Validation, logger.Api, "SendOTP failed due to invalid phone number", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid phone number", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "SendOTP internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to send OTP", err)
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for VerifyOTP", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			h.log.Warn(logger.Validation, logger.Api, "VerifyOTP failed: invalid code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			h.log.Warn(logger.Validation, logger.Api, "VerifyOTP failed: expired code or phone not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "OTP code has expired or phone not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "VerifyOTP internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for Login", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrNotVerifiedPhoneNumber):
			h.log.Warn(logger.Validation, logger.Api, "Login attempt on unverified phone number", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, "Your phone number has not verified yet", err)
		case errors.Is(err, exception.ErrInvalidCredentials):
			h.log.Warn(logger.Validation, logger.Api, "Login failed: invalid credentials", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, "Invalid phone number or password", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "Login internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for Register", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req.Phone, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserAlreadyExists):
			h.log.Warn(logger.Validation, logger.Api, "Register failed: user already exists", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "User already exists", err)
		case errors.Is(err, exception.ErrInvalidPhone):
			h.log.Warn(logger.Validation, logger.Api, "Register failed: invalid phone", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid phone number", err)
		case errors.Is(err, exception.ErrInvalidPassword):
			h.log.Warn(logger.Validation, logger.Api, "Register failed: invalid password", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid password", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "Register internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for Refresh", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidToken):
			h.log.Warn(logger.Validation, logger.Api, "RefreshToken failed: invalid token", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, "Invalid refresh token", err)
		case errors.Is(err, exception.ErrTokenExpired):
			h.log.Warn(logger.Validation, logger.Api, "RefreshToken failed: expired token", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, "Refresh token has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "RefreshToken failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "RefreshToken internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to refresh token", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", tokens)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userid, err := h.contexttokenService.GetUserIDFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "Logout failed: could not extract user ID context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	jwtData, err := h.contexttokenService.GetJWTFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "Logout failed: could not extract JWT context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetJWTDataFromContext.Error())
		return
	}

	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for Logout", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userid, &jwtData, &req); err != nil {
		if errors.Is(err, exception.ErrInvalidToken) {
			h.log.Warn(logger.Validation, logger.Api, "Logout failed: invalid token", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid token", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "Logout internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to logout", err)
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}

// VerifyPhone handles POST /auth/verify-phone
func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	userid, err := h.contexttokenService.GetUserIDFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "VerifyPhone failed: missing user token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for VerifyPhone", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyPhone(c.Request.Context(), userid, req.Phone, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			h.log.Warn(logger.Validation, logger.Api, "VerifyPhone failed: invalid OTP code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			h.log.Warn(logger.Validation, logger.Api, "VerifyPhone failed: OTP code expired", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "OTP code has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "VerifyPhone failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "VerifyPhone internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
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
		h.log.Warn(logger.Internal, logger.Api, "VerifyEmail failed: missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
		return
	}

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid request payload for VerifyEmail", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	success, err := h.authService.VerifyEmail(c.Request.Context(), &tokenData, req.Email, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrInvalidOTP):
			h.log.Warn(logger.Validation, logger.Api, "VerifyEmail failed: invalid OTP code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Invalid OTP code", err)
		case errors.Is(err, exception.ErrOTPExpired):
			h.log.Warn(logger.Validation, logger.Api, "VerifyEmail failed: OTP expired", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "OTP code has expired", err)
		case errors.Is(err, exception.ErrUserNotFound):
			h.log.Warn(logger.Validation, logger.Api, "VerifyEmail failed: user not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "User not found", err)
		default:
			h.log.Error(logger.Internal, logger.Api, "VerifyEmail internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to verify email address", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Email address verified successfully", gin.H{
		"success": success,
	})
}

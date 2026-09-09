package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type AuthHandler struct {
	authService         service_contract.AuthenticationServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewAuthHandler(authService service_contract.AuthenticationServiceInterface, contexttokenService service_contract.ContextTokenServiceInterface, log domainlogger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, contexttokenService: contexttokenService, log: log}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.authService.SendOTP(c.Request.Context(), req.Phone)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "OTP sent successfully", result)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.authService.VerifyOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Authentication successful", result)
}

func (h *AuthHandler) SendEmailOTP(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.authService.SendEmailOTP(c.Request.Context(), userID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Email OTP sent successfully", result)
}

func (h *AuthHandler) VerifyEmailOTP(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	var req dto.VerifyEmailOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	if err := h.authService.VerifyEmailOTP(c.Request.Context(), userID, req.Code); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Email verified successfully", nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Token refreshed successfully", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, err := h.contexttokenService.GetUserIDFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}
	jwtData, err := h.contexttokenService.GetJWTFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetJWTDataFromContext, h.log)
		return
	}
	var req dto.LogoutRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.authService.Logout(c.Request.Context(), userID, &jwtData, &req); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}

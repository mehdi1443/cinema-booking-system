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

func NewAuthHandler(
	authService service_contract.AuthenticationServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// POST /auth/otp/send
// -----------------------------------------------------------------------------

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if (req.Email == req.Phone && req.Email == "") || (req.Email != "" && req.Phone != "") {
		HandleBadRequest(c, "Just fill exactly one field (email or phone)")
		return
	}

	success, err := h.authService.SendOTP(c.Request.Context(), req.Phone, req.Email)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "OTP sent successfully", gin.H{
		"success": success,
	})
}

// -----------------------------------------------------------------------------
// POST /auth/otp/verify
// -----------------------------------------------------------------------------

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	success, err := h.authService.VerifyOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "OTP verified successfully", gin.H{
		"success": success,
	})
}

// -----------------------------------------------------------------------------
// POST /auth/login
// -----------------------------------------------------------------------------

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

// -----------------------------------------------------------------------------
// POST /auth/register
// -----------------------------------------------------------------------------

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req.Phone, req.Email, req.Password)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Registration successful", res)
}

// -----------------------------------------------------------------------------
// POST /auth/refresh
// -----------------------------------------------------------------------------

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", tokens)
}

// -----------------------------------------------------------------------------
// POST /auth/logout
// -----------------------------------------------------------------------------

func (h *AuthHandler) Logout(c *gin.Context) {
	userid, err := h.contexttokenService.GetUserIDFromGinContext(c)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userid, &jwtData, &req); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}

// -----------------------------------------------------------------------------
// POST /auth/verify-phone
// -----------------------------------------------------------------------------

func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	userid, err := h.contexttokenService.GetUserIDFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	success, err := h.authService.VerifyPhone(c.Request.Context(), userid, req.Phone, req.Code)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Phone number verified successfully", gin.H{
		"success": success,
	})
}

// -----------------------------------------------------------------------------
// POST /auth/verify-email
// -----------------------------------------------------------------------------

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, exception.ErrGetTokenFromContext, h.log)
		return
	}

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	success, err := h.authService.VerifyEmail(c.Request.Context(), &tokenData, req.Email, req.Code)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Email address verified successfully", gin.H{
		"success": success,
	})
}

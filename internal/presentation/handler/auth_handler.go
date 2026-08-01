package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/application/auth"
)

type AuthHandler struct {
	authService auth.MyAuthService
}

func NewAuthHandler(authService auth.MyAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) SendOTPHandler(c *gin.Context) {
	
}

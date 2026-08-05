package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/otp/send", authHandler.SendOTP)
		auth.POST("/otp/verify", authHandler.VerifyOTP)
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)

		auth.POST("/verify-phone", authHandler.VerifyPhone)
		auth.POST("/verify-email", authHandler.VerifyEmail)
	}
}

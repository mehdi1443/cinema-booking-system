package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpAuthRouter(
	router *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	captchaMiddleware *middleware.CaptchaSolving,
) {
	auth := router.Group("/auth")
	auth.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		auth.POST("/otp/send", captchaMiddleware.CaptchaSolvingMiddleware(), authHandler.SendOTP)
		auth.POST("/otp/verify", authHandler.VerifyOTP)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}
}

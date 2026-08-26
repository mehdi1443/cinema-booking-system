package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpCaptchaRouter(
	router *gin.RouterGroup,
	captchaHandler *handler.CaptchaHandler,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	appLogger domainlogger.Logger,
) {
	captcha := router.Group("/captcha")
	captcha.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		captcha.POST("/send",
			captchaHandler.SendCaptcha,
		)
		captcha.POST("/verify",
			captchaHandler.VerifyCaptcha,
		)
	}
}

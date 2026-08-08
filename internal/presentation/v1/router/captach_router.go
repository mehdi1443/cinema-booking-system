package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

func SetUpCaptchaRouter(
	router *gin.RouterGroup,
	captchaHandler *handler.CaptchaHandler,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	appLogger logger.Logger,
) {

	captcha := router.Group("/captcha")
	captcha.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		captcha.POST("/send",
			captchaHandler.SendCaptcha,
		)
		captcha.POST("/",
			captchaHandler.VerifyCaptcha,
		)

	}

}

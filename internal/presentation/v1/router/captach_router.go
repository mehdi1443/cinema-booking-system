package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpCaptchaRouter(router *gin.RouterGroup, captchaHandler *handler.CaptchaHandler) {
	captcha := router.Group("/captcha")
	{
		captcha.POST("/send",
			captchaHandler.SendCaptcha,
		)
		captcha.POST("/verify",
			captchaHandler.VerifyCaptcha,
		)
	}
}

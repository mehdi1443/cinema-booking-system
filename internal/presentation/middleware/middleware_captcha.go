package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

const CtxCaptchaValidationKey = "CAPTCHA_VALIDATION"
const CtxCaptchaIDKey = "X-Captcha-ID"
const CtxCaptchaCodeKey = "X-Captcha-Code"

type CaptchaSolving struct {
	captchaService service_contract.CaptchaServiceInterface
	debugConfig    *bootstrap.DebugModeOptions
	log            logger.Logger
}

func NewCaptchaSolving(
	captchaService service_contract.CaptchaServiceInterface,
	debugConfig *bootstrap.DebugModeOptions,
	log logger.Logger,
) *CaptchaSolving {
	return &CaptchaSolving{
		captchaService: captchaService,
		debugConfig:    debugConfig,
		log:            log,
	}
}

func (cs *CaptchaSolving) CaptchaSolvingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cs.debugConfig.Enable {
			c.Set(CtxCaptchaValidationKey, true)
			c.Next()
			return
		}

		captchaID := c.GetHeader(CtxCaptchaIDKey)
		captchaCode := c.GetHeader(CtxCaptchaCodeKey)

		if captchaID == "" || captchaCode == "" {
			cs.log.Warn(logger.Validation, logger.Api, "Captcha headers missing", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: "X-Captcha-ID or X-Captcha-Code header is empty",
			})
			response.Error(c, http.StatusBadRequest, "captcha_headers_required", nil)
			c.Abort()
			return
		}

		isValid, err := cs.captchaService.Verify(c.Request.Context(), captchaID, captchaCode)
		if err != nil {
			if errors.Is(err, exception.ErrCaptchaNotFound) {
				cs.log.Warn(logger.Validation, logger.Api, "VerifyCaptcha failed: captcha not found or expired", map[logger.ExtraKey]interface{}{
					logger.ErrorMessage: err.Error(),
				})
				response.Error(c, http.StatusNotFound, "Captcha not found or expired", err)
				c.Abort()
				return
			}

			cs.log.Error(logger.Internal, logger.Api, "VerifyCaptcha internal failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Failed to verify captcha", err)
			c.Abort()
			return
		}

		if isValid == nil || !*isValid {
			cs.log.Warn(logger.Validation, logger.Api, "VerifyCaptcha failed: invalid captcha code", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: "code mismatch",
			})
			response.Error(c, http.StatusBadRequest, "Invalid captcha code", nil)
			c.Abort()
			return
		}

		c.Set(CtxCaptchaValidationKey, true)
		c.Next()
	}
}

package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

const CtxCaptchaValidationKey = "CAPTCHA_VALIDATION"
const CtxCaptchaIDKey = "X-Captcha-ID"
const CtxCaptchaCodeKey = "X-Captcha-Code"

type CaptchaSolving struct {
	captchaService service_contract.CaptchaServiceInterface
	debugConfig    *bootstrap.DebugModeOptions
	log            domainlogger.Logger
}

func NewCaptchaSolving(
	captchaService service_contract.CaptchaServiceInterface,
	debugConfig *bootstrap.DebugModeOptions,
	log domainlogger.Logger,
) *CaptchaSolving {
	return &CaptchaSolving{
		captchaService: captchaService,
		debugConfig:    debugConfig,
		log:            log,
	}
}

func (cs *CaptchaSolving) CaptchaSolvingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if cs.debugConfig.Enable {
			c.Set(CtxCaptchaValidationKey, true)
			c.Next()
			return
		}

		captchaID := c.GetHeader(CtxCaptchaIDKey)
		captchaCode := c.GetHeader(CtxCaptchaCodeKey)

		if captchaID == "" || captchaCode == "" {
			cs.log.Warn(ctx, "Captcha headers missing", domainlogger.String("details", "X-Captcha-ID or X-Captcha-Code header is empty"))
			response.ErrorWithCode(c, http.StatusBadRequest, string(exception.CodeInvalidInput), "Captcha headers are required")
			c.Abort()
			return
		}

		isValid, err := cs.captchaService.Verify(ctx, captchaID, captchaCode)
		if err != nil {
			if errors.Is(err, exception.ErrCaptchaNotFound) {
				cs.log.Warn(ctx, "VerifyCaptcha failed: captcha not found or expired", domainlogger.Error(err))
				response.ErrorWithCode(c, http.StatusNotFound, string(exception.CodeCaptchaNotFound), "Captcha not found or expired")
				c.Abort()
				return
			}

			cs.log.Error(ctx, "VerifyCaptcha internal failure", domainlogger.Error(err))
			response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "Failed to verify captcha")
			c.Abort()
			return
		}

		if isValid == nil || !*isValid {
			cs.log.Warn(ctx, "VerifyCaptcha failed: invalid captcha code", domainlogger.String("details", "code mismatch"))
			response.ErrorWithCode(c, http.StatusBadRequest, string(exception.CodeInvalidInput), "Invalid captcha code")
			c.Abort()
			return
		}

		c.Set(CtxCaptchaValidationKey, true)
		c.Next()
	}
}

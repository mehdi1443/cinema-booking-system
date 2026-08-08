package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type CaptchaHandler struct {
	captchaService      service_contract.CaptchaServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewCaptchaHandler(
	captchaService service_contract.CaptchaServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *CaptchaHandler {
	return &CaptchaHandler{
		captchaService:      captchaService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// SendCaptcha handles POST /captcha/send
func (h *CaptchaHandler) SendCaptcha(c *gin.Context) {
	captcha, err := h.captchaService.GenerateAndSave(c.Request.Context())
	if err != nil {
		h.log.Error(logger.Internal, logger.Api, "SendCaptcha internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to generate captcha", err)
		return
	}

	response.Success(c, http.StatusOK, "Captcha generated successfully", captcha)
}

// VerifyCaptcha handles POST /captcha/verify
func (h *CaptchaHandler) VerifyCaptcha(c *gin.Context) {
	var req dto.VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for VerifyCaptcha", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	isValid, err := h.captchaService.Verify(c.Request.Context(), req.CaptchaID, req.Code)
	if err != nil {
		if errors.Is(err, exception.ErrCaptchaNotFound) {
			h.log.Warn(logger.Validation, logger.Api, "VerifyCaptcha failed: captcha not found or expired", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Captcha not found or expired", err)
			return
		}

		h.log.Error(logger.Internal, logger.Api, "VerifyCaptcha internal failure", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to verify captcha", err)
		return
	}

	if isValid == nil || !*isValid {
		h.log.Warn(logger.Validation, logger.Api, "VerifyCaptcha failed: invalid captcha code", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: "code mismatch",
		})
		response.Error(c, http.StatusBadRequest, "Invalid captcha code", nil)
		return
	}

	response.Success(c, http.StatusOK, "Captcha verified successfully", dto.VerifyCaptchaResponse{
		Valid: true,
	})
}

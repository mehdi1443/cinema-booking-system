package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type CaptchaHandler struct {
	captchaService      service_contract.CaptchaServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewCaptchaHandler(
	captchaService service_contract.CaptchaServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *CaptchaHandler {
	return &CaptchaHandler{
		captchaService:      captchaService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// POST /captcha/send
// -----------------------------------------------------------------------------

func (h *CaptchaHandler) SendCaptcha(c *gin.Context) {
	captcha, err := h.captchaService.GenerateAndSave(c.Request.Context())
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Captcha generated successfully", captcha)
}

// -----------------------------------------------------------------------------
// POST /captcha/verify
// -----------------------------------------------------------------------------

func (h *CaptchaHandler) VerifyCaptcha(c *gin.Context) {
	var req dto.VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	isValid, err := h.captchaService.Verify(c.Request.Context(), req.CaptchaID, req.Code)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	if isValid == nil || !*isValid {
		HandleBadRequest(c, "Invalid captcha code")
		return
	}

	response.Success(c, http.StatusOK, "Captcha verified successfully", dto.VerifyCaptchaResponse{
		Valid: true,
	})
}

package dto

import "time"

type GenerateCaptcha struct {
	CaptchaId string    `json:"captcha_id"`
	Code      string    `json:"captcha_code"`
	ExpireAt  time.Time `json:"expired_at"`
}

type VerifyCaptchaRequest struct {
	CaptchaID string `json:"captcha_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type VerifyCaptchaResponse struct {
	Valid bool `json:"valid"`
}

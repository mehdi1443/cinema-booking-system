package dto

import "time"

type GenerateCaptcha struct {
	CaptchaId string    `json:"captcha_id"`
	ImageCode string    `json:"captcha_img_code_base64,omitempty"`
	AudioCode string    `json:"captcha_audio_code_base64,omitempty"`
	ExpireAt  time.Time `json:"expired_at"`
}

type VerifyCaptchaRequest struct {
	CaptchaID string `json:"captcha_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type VerifyCaptchaResponse struct {
	Valid bool `json:"valid"`
}

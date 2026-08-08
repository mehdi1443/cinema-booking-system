package captcha

import (
	"fmt"

	"github.com/mojocn/base64Captcha"
)

type CaptchaMan interface {
	Generate() (cid string, b64 string, answer string, err error)
}

type captchaMan struct {
	captcha *base64Captcha.Captcha
}

// captchaType == "string" | "math" | "audio"
func NewCaptchaMan(captchaType string) (CaptchaMan, error) {
	dummyStore := &dummyStore{}
	var c *base64Captcha.Captcha

	switch captchaType {
	case "string":
		cfg := suggestDriverString
		driver := base64Captcha.NewDriverString(cfg.Height, cfg.Width, cfg.NoiseCount, cfg.ShowLineOptions, cfg.Length, cfg.Source, cfg.BgColor, cfg.FontsStorage, cfg.Fonts)
		c = base64Captcha.NewCaptcha(driver, dummyStore)

	case "math":
		cfg := suggestDriverMath
		driver := base64Captcha.NewDriverMath(cfg.Height, cfg.Width, cfg.NoiseCount, cfg.ShowLineOptions, cfg.BgColor, cfg.FontsStorage, cfg.Fonts)
		c = base64Captcha.NewCaptcha(driver, dummyStore)

	case "audio":
		cfg := suggestDriverAudio
		driver := base64Captcha.NewDriverAudio(cfg.Length, cfg.Language)
		c = base64Captcha.NewCaptcha(driver, dummyStore)

	default:
		return nil, fmt.Errorf("invalid captcha-type")
	}

	return &captchaMan{
		captcha: c,
	}, nil
}

func (cm *captchaMan) Generate() (cid string, b64 string, answer string, err error) {
	return cm.captcha.Generate()
}

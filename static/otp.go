package static

import (
	"bytes"
)

func (sf *StaticFiles) ExecuteOTPTemplate(otpCode string, expireMinutes int) (*bytes.Buffer, error) {
	var htmlBuffer bytes.Buffer
	if err := sf.otpHTMLTemplate.Execute(&htmlBuffer, map[string]any{
		"OTPCode":          otpCode,
		"OTPExpireMinutes": expireMinutes,
	}); err != nil {
		return nil, err
	}
	return &htmlBuffer, nil
}

package validation

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	otpCodeRegex *regexp.Regexp
)

func ValidateOTPCode(codeLen int) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		otpCodeRegex = regexp.MustCompile(fmt.Sprintf(`^\d{%d}$`, codeLen))
		code := fl.Field().String()
		return otpCodeRegex.MatchString(code)

	}
}

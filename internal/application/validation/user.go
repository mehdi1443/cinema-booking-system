package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	// Matches 09123456789, +989123456789, or 9123456789
	iranianPhoneRegex = regexp.MustCompile(`^(\+98|0)?9\d{9}$`)
)

// validateIranianPhone validates Iranian mobile number formats.
func ValidateIranianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return iranianPhoneRegex.MatchString(phone)
}

// validatePassword enforces min 8 characters, at least 1 uppercase letter,
// 1 lowercase letter, and 1 digit.
func ValidatePassword(debugMode bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if debugMode == true {
			if password != "" {
				return true
			}
			return false
		}
		if len(password) < 8 {
			return false
		}

		var hasUpper, hasLower, hasDigit bool
		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			}
		}

		return hasUpper && hasLower && hasDigit
	}
}

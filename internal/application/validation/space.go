package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateNotBlank(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return strings.TrimSpace(str) != ""
}

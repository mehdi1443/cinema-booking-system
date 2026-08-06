package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	// حروف الفبای فارسی: (پ، چ، ژ، گ، ک، ی و...)
	// حروف عربی رایج در کیبوردهای مختلف: (ك، ي، ة، أ، إ، ؤ، ء و...)
	// اعداد: اعداد فارسی (۰-۹) و اعداد عربی (٠-٩)
	// فاصله‌ها: فاصله معمولی (Space) و نیم‌فاصله (\u200C - ZWNJ)
	// علائم نگارشی فارسی/عربی: (، ؛ ؟ % و...)
	// رسما میشه فقط متن فارسی
	roleTitleRegex = regexp.MustCompile(
		`^[ \t\u0621-\u0628\u062A-\u063A\u0641-\u0642\u0644-\u0648\u064E-\u0651\u0655\u067E\u0686\u0698\u06A9\u06AF\u06BE\u06CC\u0629\u0643\u0649-\u064B\u064D\u06D5\u06F0-\u06F9\u0660-\u0669\u200C\u200D\u060C\u061B\u061F\u0640\u066A\u066B\u066C]{2,}$`)
)

func ValidateRoleTitle(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return roleTitleRegex.MatchString(title)
}

package exception

import "errors"

var (
	ErrCaptchaNotFound         = errors.New("captcha not found or expired")
	ErrCaptchaExists           = errors.New("captcha already exists")
	ErrCaptchaGenerationFailed = errors.New("failed to generate captcha")

	// Auth & Token Errors
	ErrRefreshTokenNotFound  = errors.New("refresh token not found in cache")
	ErrTokenExpired          = errors.New("token has expired")
	ErrInvalidToken          = errors.New("invalid or malformed token")
	ErrGetTokenFromContext   = errors.New("failed to get token from context")
	ErrGetJWTDataFromContext = errors.New("failed to get JWT data from context")

	// Access & Validation Errors
	ErrNoteAccessible = errors.New("insufficient access permissions")
	ErrInvalidFilter  = errors.New("invalid search filter criteria")

	// User Credentials & OTP
	ErrInvalidPassword        = errors.New("invalid user password")
	ErrInvalidPhone           = errors.New("invalid user phone number")
	ErrInvalidCredentials     = errors.New("invalid user credentials")
	ErrNotVerifiedPhoneNumber = errors.New("phone number is not verified")
	ErrInvalidOTP             = errors.New("invalid otp code")
	ErrOTPExpired             = errors.New("otp code has expired")
	ErrOtpKeyNotFound         = errors.New("otp key not found in cache")

	// Infrastructure & Generic Errors
	ErrDatabase     = errors.New("database operation failed")
	ErrCache        = errors.New("cache operation failed")
	ErrDuplicateKey = errors.New("record already exists")

	// User Errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUsersNotFound     = errors.New("no users found matching criteria")

	// RBAC Errors
	ErrRoleNotFound            = errors.New("role not found")
	ErrRoleAlreadyExists       = errors.New("role already exists")
	ErrRolesNotFound           = errors.New("no roles found matching criteria")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrPermissionsNotFound     = errors.New("no permissions found matching criteria")
	ErrUserPermissionNotFound  = errors.New("user permission not found in cache")
	ErrAssociationFailed       = errors.New("failed to process association operation")
	ErrRolePermissionNotFound  = errors.New("role-permission association not found")
	ErrUserRoleNotFound        = errors.New("user-role association not found")

	// Wallet & Balance Errors
	ErrInvalidAmount       = errors.New("invalid transaction amount")
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrReachManBalance     = errors.New("are you reach or something?!")

	// Cinema Domain Errors
	ErrMovieNotFound      = errors.New("movie not found")
	ErrMovieAlreadyExists = errors.New("movie already exists")

	ErrGenreNotFound      = errors.New("genre not found")
	ErrGenreAlreadyExists = errors.New("genre already exists")

	ErrPersonNotFound      = errors.New("person not found")
	ErrPersonAlreadyExists = errors.New("person already exists")

	ErrInvalidFileType = errors.New("invalid file type provided")
)

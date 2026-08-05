package exception

import "errors"

var (
	ErrRefreshTokenNotFound   = errors.New("refresh token not found in redis")
	ErrUserPermissionNotFound = errors.New("user permission not found in redis")
	ErrOtpKeyNotFound         = errors.New("otp key not found in redis")
	ErrTokenExpired           = errors.New("token has expired")
	ErrInvalidToken           = errors.New("invalid or malformed token")
	ErrGetTokenFromContext    = errors.New("fail get keys from context")

	ErrNoteAccessible = errors.New("you have not enough access")
	ErrInvalidFilter  = errors.New("invalid search filter criteria")

	ErrInvalidPassword    = errors.New("invalid user password")
	ErrInvalidPhone       = errors.New("invalid user phone")
	ErrInvalidCredentials = errors.New("invalid user credentials")
	ErrInvalidOTP         = errors.New("invalid otp code")
	ErrOTPExpired         = errors.New("otp has expired")

	//================================================
	// ErrDatabase
	//================================================

	ErrDatabase     = errors.New("database operation failed")
	ErrDuplicateKey = errors.New("record already exists")

	//================================================
	// ErrUserNotFound
	//================================================

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")

	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")

	//================================================
	// ErrInvalidAmount
	//================================================

	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrReachManBalance     = errors.New("are you reach or something?!")

	//================================================
	// Cinema Errors
	//================================================

	ErrMovieNotFound  = errors.New("movie not found")
	ErrGenreNotFound  = errors.New("genre not found")
	ErrPersonNotFound = errors.New("person not found")
)

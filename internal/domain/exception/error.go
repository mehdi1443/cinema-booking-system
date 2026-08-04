package exception

import "errors"

var (
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

	//================================================
	// Cinema Errors
	//================================================

	ErrMovieNotFound  = errors.New("movie not found")
	ErrGenreNotFound  = errors.New("genre not found")
	ErrPersonNotFound = errors.New("person not found")
)

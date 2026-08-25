package exception

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeInvalidInput  ErrorCode = "INVALID_INPUT"
	CodeInvalidFilter ErrorCode = "INVALID_FILTER"

	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"

	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeUserNotFound   ErrorCode = "USER_NOT_FOUND"
	CodeMovieNotFound  ErrorCode = "MOVIE_NOT_FOUND"
	CodeGenreNotFound  ErrorCode = "GENRE_NOT_FOUND"
	CodePersonNotFound ErrorCode = "PERSON_NOT_FOUND"

	CodeConflict          ErrorCode = "CONFLICT"
	CodeDuplicateKey      ErrorCode = "DUPLICATE_KEY"
	CodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"

	CodeInvalidPassword ErrorCode = "INVALID_PASSWORD"

	CodeDatabase ErrorCode = "DATABASE_ERROR"
	CodeInternal ErrorCode = "INTERNAL_ERROR"

	CodeTokenContext ErrorCode = "TOKEN_CONTEXT_INVALID"

	CodeMovieAlreadyExists ErrorCode = "MOVIE_ALREADY_EXISTS"
	CodeInvalidFileType    ErrorCode = "INVALID_FILE_TYPE"
	CodeFileSizeTooLarge   ErrorCode = "FILE_SIZE_TOO_LARGE"

	CodeGenreAlreadyExists  ErrorCode = "GENRE_ALREADY_EXISTS"
	CodePersonAlreadyExists ErrorCode = "PERSON_ALREADY_EXISTS"
	CodeResourceInUse       ErrorCode = "RESOURCE_IN_USE"
)

type ErrorKind string

const (
	KindValidation ErrorKind = "VALIDATION"
	KindAuth       ErrorKind = "AUTH"
	KindForbidden  ErrorKind = "FORBIDDEN"
	KindNotFound   ErrorKind = "NOT_FOUND"
	KindConflict   ErrorKind = "CONFLICT"
	KindDatabase   ErrorKind = "DATABASE"
	KindInternal   ErrorKind = "INTERNAL"
)

type AppError struct {
	Code    ErrorCode
	Kind    ErrorKind
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}

	return e.Code == t.Code
}

func New(
	code ErrorCode,
	kind ErrorKind,
	message string,
) *AppError {
	return &AppError{
		Code:    code,
		Kind:    kind,
		Message: message,
	}
}

func Wrap(
	target error,
	message string,
	cause error,
) error {
	if cause == nil {
		return target
	}

	var appErr *AppError

	if errors.As(target, &appErr) {
		return &AppError{
			Code:    appErr.Code,
			Kind:    appErr.Kind,
			Message: message,
			Cause:   cause,
		}
	}

	return &AppError{
		Code:    CodeInternal,
		Kind:    KindInternal,
		Message: message,
		Cause:   cause,
	}
}

// -----------------------------------------------------------------------------
// Authentication / Authorization
// -----------------------------------------------------------------------------

var ErrGetTokenFromContext = &AppError{
	Code:    CodeTokenContext,
	Kind:    KindAuth,
	Message: "invalid authentication context",
}

var ErrNoteAccessible = &AppError{
	Code:    CodeForbidden,
	Kind:    KindForbidden,
	Message: "access denied",
}

// -----------------------------------------------------------------------------
// Users
// -----------------------------------------------------------------------------

var ErrUserNotFound = &AppError{
	Code:    CodeUserNotFound,
	Kind:    KindNotFound,
	Message: "user not found",
}

var ErrUserAlreadyExists = &AppError{
	Code:    CodeUserAlreadyExists,
	Kind:    KindConflict,
	Message: "user already exists",
}

var ErrInvalidPassword = &AppError{
	Code:    CodeInvalidPassword,
	Kind:    KindValidation,
	Message: "invalid password",
}

// -----------------------------------------------------------------------------
// Movies
// -----------------------------------------------------------------------------

var ErrMovieNotFound = &AppError{
	Code:    CodeMovieNotFound,
	Kind:    KindNotFound,
	Message: "movie not found",
}

var ErrGenreNotFound = &AppError{
	Code:    CodeGenreNotFound,
	Kind:    KindNotFound,
	Message: "genre not found",
}

var ErrPersonNotFound = &AppError{
	Code:    CodePersonNotFound,
	Kind:    KindNotFound,
	Message: "person not found",
}

// -----------------------------------------------------------------------------
// Generic
// -----------------------------------------------------------------------------

var ErrDuplicateKey = &AppError{
	Code:    CodeDuplicateKey,
	Kind:    KindConflict,
	Message: "duplicate key",
}

var ErrInvalidFilter = &AppError{
	Code:    CodeInvalidFilter,
	Kind:    KindValidation,
	Message: "invalid filter",
}

var ErrDatabase = &AppError{
	Code:    CodeDatabase,
	Kind:    KindDatabase,
	Message: "database operation failed",
}

var ErrInternal = &AppError{
	Code:    CodeInternal,
	Kind:    KindInternal,
	Message: "internal server error",
}

var ErrGenreAlreadyExists = &AppError{
	Code:    CodeGenreAlreadyExists,
	Kind:    KindConflict,
	Message: "genre already exists",
}

var ErrPersonAlreadyExists = &AppError{
	Code:    CodePersonAlreadyExists,
	Kind:    KindConflict,
	Message: "person already exists",
}

var ErrResourceInUse = &AppError{
	Code:    CodeResourceInUse,
	Kind:    KindConflict,
	Message: "resource is in use and cannot be modified or deleted",
}

var ErrInvalidInput = &AppError{
	Code:    CodeInvalidInput,
	Kind:    KindValidation,
	Message: "invalid input",
}

var ErrMovieAlreadyExists = &AppError{
	Code:    CodeMovieAlreadyExists,
	Kind:    KindConflict,
	Message: "movie already exists",
}

var ErrInvalidFileType = &AppError{
	Code:    CodeInvalidFileType,
	Kind:    KindValidation,
	Message: "invalid file type",
}

var ErrFileSizeTooLarge = &AppError{
	Code:    CodeFileSizeTooLarge,
	Kind:    KindValidation,
	Message: "file size is too large",
}

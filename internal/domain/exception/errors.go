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

	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	CodeMovieNotFound    ErrorCode = "MOVIE_NOT_FOUND"
	CodeGenreNotFound    ErrorCode = "GENRE_NOT_FOUND"
	CodePersonNotFound   ErrorCode = "PERSON_NOT_FOUND"
	CodeCinemaNotFound   ErrorCode = "CINEMA_NOT_FOUND"
	CodeFacilityNotFound ErrorCode = "FACILITY_NOT_FOUND"
	CodeCityNotFound     ErrorCode = "CITY_NOT_FOUND"
	CodeProvinceNotFound ErrorCode = "PROVINCE_NOT_FOUND"
	CodeCommentNotFound  ErrorCode = "COMMENT_NOT_FOUND"
	CodeFavoriteNotFound ErrorCode = "FAVORITE_NOT_FOUND"
	CodeRatingNotFound   ErrorCode = "RATING_NOT_FOUND"

	CodeConflict              ErrorCode = "CONFLICT"
	CodeDuplicateKey          ErrorCode = "DUPLICATE_KEY"
	CodeUserAlreadyExists     ErrorCode = "USER_ALREADY_EXISTS"
	CodeMovieAlreadyExists    ErrorCode = "MOVIE_ALREADY_EXISTS"
	CodeGenreAlreadyExists    ErrorCode = "GENRE_ALREADY_EXISTS"
	CodePersonAlreadyExists   ErrorCode = "PERSON_ALREADY_EXISTS"
	CodeCinemaAlreadyExists   ErrorCode = "CINEMA_ALREADY_EXISTS"
	CodeFacilityAlreadyExists ErrorCode = "FACILITY_ALREADY_EXISTS"
	CodeCityAlreadyExists     ErrorCode = "CITY_ALREADY_EXISTS"
	CodeProvinceAlreadyExists ErrorCode = "PROVINCE_ALREADY_EXISTS"
	CodeCommentAlreadyExists  ErrorCode = "COMMENT_ALREADY_EXISTS"
	CodeFavoriteAlreadyExists ErrorCode = "FAVORITE_ALREADY_EXISTS"
	CodeRatingAlreadyExists   ErrorCode = "RATING_ALREADY_EXISTS"
	CodeResourceInUse         ErrorCode = "RESOURCE_IN_USE"

	CodeInvalidPassword  ErrorCode = "INVALID_PASSWORD"
	CodeDatabase         ErrorCode = "DATABASE_ERROR"
	CodeInternal         ErrorCode = "INTERNAL_ERROR"
	CodeTokenContext     ErrorCode = "TOKEN_CONTEXT_INVALID"
	CodeInvalidFileType  ErrorCode = "INVALID_FILE_TYPE"
	CodeFileSizeTooLarge ErrorCode = "FILE_SIZE_TOO_LARGE"
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
// Movies & Entities
// -----------------------------------------------------------------------------

var ErrMovieNotFound = &AppError{
	Code:    CodeMovieNotFound,
	Kind:    KindNotFound,
	Message: "movie not found",
}

var ErrMovieAlreadyExists = &AppError{
	Code:    CodeMovieAlreadyExists,
	Kind:    KindConflict,
	Message: "movie already exists",
}

var ErrGenreNotFound = &AppError{
	Code:    CodeGenreNotFound,
	Kind:    KindNotFound,
	Message: "genre not found",
}

var ErrGenreAlreadyExists = &AppError{
	Code:    CodeGenreAlreadyExists,
	Kind:    KindConflict,
	Message: "genre already exists",
}

var ErrPersonNotFound = &AppError{
	Code:    CodePersonNotFound,
	Kind:    KindNotFound,
	Message: "person not found",
}

var ErrPersonAlreadyExists = &AppError{
	Code:    CodePersonAlreadyExists,
	Kind:    KindConflict,
	Message: "person already exists",
}

// -----------------------------------------------------------------------------
// Cinema & Infrastructure
// -----------------------------------------------------------------------------

var ErrCinemaNotFound = &AppError{
	Code:    CodeCinemaNotFound,
	Kind:    KindNotFound,
	Message: "cinema not found",
}

var ErrCinemaAlreadyExists = &AppError{
	Code:    CodeCinemaAlreadyExists,
	Kind:    KindConflict,
	Message: "cinema already exists",
}

var ErrFacilityNotFound = &AppError{
	Code:    CodeFacilityNotFound,
	Kind:    KindNotFound,
	Message: "facility not found",
}

var ErrFacilityAlreadyExists = &AppError{
	Code:    CodeFacilityAlreadyExists,
	Kind:    KindConflict,
	Message: "facility already exists",
}

var ErrCityNotFound = &AppError{
	Code:    CodeCityNotFound,
	Kind:    KindNotFound,
	Message: "city not found",
}

var ErrCityAlreadyExists = &AppError{
	Code:    CodeCityAlreadyExists,
	Kind:    KindConflict,
	Message: "city already exists",
}

var ErrProvinceNotFound = &AppError{
	Code:    CodeProvinceNotFound,
	Kind:    KindNotFound,
	Message: "province not found",
}

var ErrProvinceAlreadyExists = &AppError{
	Code:    CodeProvinceAlreadyExists,
	Kind:    KindConflict,
	Message: "province already exists",
}

// -----------------------------------------------------------------------------
// Comments, Favorites & Ratings
// -----------------------------------------------------------------------------

var ErrCommentNotFound = &AppError{
	Code:    CodeCommentNotFound,
	Kind:    KindNotFound,
	Message: "comment not found",
}

var ErrCommentAlreadyExists = &AppError{
	Code:    CodeCommentAlreadyExists,
	Kind:    KindConflict,
	Message: "comment already exists",
}

var ErrFavoriteNotFound = &AppError{
	Code:    CodeFavoriteNotFound,
	Kind:    KindNotFound,
	Message: "favorite not found",
}

var ErrFavoriteAlreadyExists = &AppError{
	Code:    CodeFavoriteAlreadyExists,
	Kind:    KindConflict,
	Message: "favorite already exists",
}

var ErrRatingNotFound = &AppError{
	Code:    CodeRatingNotFound,
	Kind:    KindNotFound,
	Message: "rating not found",
}

var ErrRatingAlreadyExists = &AppError{
	Code:    CodeRatingAlreadyExists,
	Kind:    KindConflict,
	Message: "rating already exists",
}

// -----------------------------------------------------------------------------
// Generic & System
// -----------------------------------------------------------------------------

var ErrDuplicateKey = &AppError{
	Code:    CodeDuplicateKey,
	Kind:    KindConflict,
	Message: "duplicate key",
}

var ErrResourceInUse = &AppError{
	Code:    CodeResourceInUse,
	Kind:    KindConflict,
	Message: "resource is in use and cannot be modified or deleted",
}

var ErrInvalidFilter = &AppError{
	Code:    CodeInvalidFilter,
	Kind:    KindValidation,
	Message: "invalid filter",
}

var ErrInvalidInput = &AppError{
	Code:    CodeInvalidInput,
	Kind:    KindValidation,
	Message: "invalid input",
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

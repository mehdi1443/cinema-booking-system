package exception

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	// Validation & General
	CodeInvalidInput   ErrorCode = "INVALID_INPUT"
	CodeInvalidFilter  ErrorCode = "INVALID_FILTER"
	CodeInvalidStatus  ErrorCode = "INVALID_STATUS"
	CodeRecordNotFound ErrorCode = "RECORD_NOT_FOUND"

	// Auth & Permissions
	CodeUnauthorized          ErrorCode = "UNAUTHORIZED"
	CodeForbidden             ErrorCode = "FORBIDDEN"
	CodeTokenContext          ErrorCode = "TOKEN_CONTEXT_INVALID"
	CodeInvalidToken          ErrorCode = "INVALID_TOKEN"
	CodeTokenExpired          ErrorCode = "TOKEN_EXPIRED"
	CodeRefreshTokenNotFound  ErrorCode = "REFRESH_TOKEN_NOT_FOUND"
	CodeGetJWTDataFromContext ErrorCode = "JWT_DATA_NOT_FOUND"
	CodeInvalidCredentials    ErrorCode = "INVALID_CREDENTIALS"

	// Users & Verification
	CodeUserNotFound           ErrorCode = "USER_NOT_FOUND"
	CodeUsersNotFound          ErrorCode = "USERS_NOT_FOUND"
	CodeUserAlreadyExists      ErrorCode = "USER_ALREADY_EXISTS"
	CodeInvalidPassword        ErrorCode = "INVALID_PASSWORD"
	CodeInvalidPhone           ErrorCode = "INVALID_PHONE"
	CodeNotVerifiedPhoneNumber ErrorCode = "PHONE_NOT_VERIFIED"
	CodeInvalidOTP             ErrorCode = "INVALID_OTP"
	CodeOTPExpired             ErrorCode = "OTP_EXPIRED"
	CodeOtpKeyNotFound         ErrorCode = "OTP_KEY_NOT_FOUND"

	// Captcha
	CodeCaptchaNotFound         ErrorCode = "CAPTCHA_NOT_FOUND"
	CodeCaptchaAlreadyExists    ErrorCode = "CAPTCHA_ALREADY_EXISTS"
	CodeCaptchaGenerationFailed ErrorCode = "CAPTCHA_GENERATION_FAILED"

	// RBAC
	CodeRoleNotFound            ErrorCode = "ROLE_NOT_FOUND"
	CodeRolesNotFound           ErrorCode = "ROLES_NOT_FOUND"
	CodeRoleAlreadyExists       ErrorCode = "ROLE_ALREADY_EXISTS"
	CodePermissionNotFound      ErrorCode = "PERMISSION_NOT_FOUND"
	CodePermissionsNotFound     ErrorCode = "PERMISSIONS_NOT_FOUND"
	CodePermissionAlreadyExists ErrorCode = "PERMISSION_ALREADY_EXISTS"
	CodeUserPermissionNotFound  ErrorCode = "USER_PERMISSION_NOT_FOUND"
	CodeRolePermissionNotFound  ErrorCode = "ROLE_PERMISSION_NOT_FOUND"
	CodeUserRoleNotFound        ErrorCode = "USER_ROLE_NOT_FOUND"
	CodeAssociationFailed       ErrorCode = "ASSOCIATION_FAILED"

	// Wallet & Transactions
	CodeInvalidAmount       ErrorCode = "INVALID_AMOUNT"
	CodeInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"
	CodeReachMaxBalance     ErrorCode = "REACH_MAX_BALANCE"

	// Movie Domain
	CodeMovieNotFound       ErrorCode = "MOVIE_NOT_FOUND"
	CodeMovieAlreadyExists  ErrorCode = "MOVIE_ALREADY_EXISTS"
	CodeGenreNotFound       ErrorCode = "GENRE_NOT_FOUND"
	CodeGenreAlreadyExists  ErrorCode = "GENRE_ALREADY_EXISTS"
	CodePersonNotFound      ErrorCode = "PERSON_NOT_FOUND"
	CodePersonAlreadyExists ErrorCode = "PERSON_ALREADY_EXISTS"

	// Cinema Domain
	CodeCinemaNotFound        ErrorCode = "CINEMA_NOT_FOUND"
	CodeCinemaAlreadyExists   ErrorCode = "CINEMA_ALREADY_EXISTS"
	CodeFacilityNotFound      ErrorCode = "FACILITY_NOT_FOUND"
	CodeFacilityAlreadyExists ErrorCode = "FACILITY_ALREADY_EXISTS"

	// Location Domain
	CodeCityNotFound          ErrorCode = "CITY_NOT_FOUND"
	CodeCityAlreadyExists     ErrorCode = "CITY_ALREADY_EXISTS"
	CodeProvinceNotFound      ErrorCode = "PROVINCE_NOT_FOUND"
	CodeProvinceAlreadyExists ErrorCode = "PROVINCE_ALREADY_EXISTS"

	// Interaction Domain
	CodeCommentNotFound       ErrorCode = "COMMENT_NOT_FOUND"
	CodeCommentAlreadyExists  ErrorCode = "COMMENT_ALREADY_EXISTS"
	CodeFavoriteNotFound      ErrorCode = "FAVORITE_NOT_FOUND"
	CodeFavoriteAlreadyExists ErrorCode = "FAVORITE_ALREADY_EXISTS"
	CodeRatingNotFound        ErrorCode = "RATING_NOT_FOUND"
	CodeRatingAlreadyExists   ErrorCode = "RATING_ALREADY_EXISTS"

	// Infrastructure & System
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeConflict         ErrorCode = "CONFLICT"
	CodeDuplicateKey     ErrorCode = "DUPLICATE_KEY"
	CodeResourceInUse    ErrorCode = "RESOURCE_IN_USE"
	CodeDatabase         ErrorCode = "DATABASE_ERROR"
	CodeCache            ErrorCode = "CACHE_ERROR"
	CodeInternal         ErrorCode = "INTERNAL_ERROR"
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
// General & Input Validation
// -----------------------------------------------------------------------------

var ErrRecordNotFound = &AppError{Code: CodeRecordNotFound, Kind: KindNotFound, Message: "record not found"}
var ErrInvalidInput = &AppError{Code: CodeInvalidInput, Kind: KindValidation, Message: "invalid input"}
var ErrInvalidFilter = &AppError{Code: CodeInvalidFilter, Kind: KindValidation, Message: "invalid search filter criteria"}
var ErrInvalidStatus = &AppError{Code: CodeInvalidStatus, Kind: KindValidation, Message: "invalid cinema status"}

// -----------------------------------------------------------------------------
// Authentication / Authorization
// -----------------------------------------------------------------------------

var ErrUnauthorized = &AppError{Code: CodeUnauthorized, Kind: KindAuth, Message: "unauthorized"}
var ErrNoteAccessible = &AppError{Code: CodeForbidden, Kind: KindForbidden, Message: "insufficient access permissions"}
var ErrInvalidCredentials = &AppError{Code: CodeInvalidCredentials, Kind: KindAuth, Message: "invalid user credentials"}

// -----------------------------------------------------------------------------
// Token & Session
// -----------------------------------------------------------------------------

var ErrGetTokenFromContext = &AppError{Code: CodeTokenContext, Kind: KindAuth, Message: "failed to get token from context"}
var ErrGetJWTDataFromContext = &AppError{Code: CodeGetJWTDataFromContext, Kind: KindAuth, Message: "failed to get JWT data from context"}
var ErrInvalidToken = &AppError{Code: CodeInvalidToken, Kind: KindAuth, Message: "invalid or malformed token"}
var ErrTokenExpired = &AppError{Code: CodeTokenExpired, Kind: KindAuth, Message: "token has expired"}
var ErrRefreshTokenNotFound = &AppError{Code: CodeRefreshTokenNotFound, Kind: KindNotFound, Message: "refresh token not found in cache"}

// -----------------------------------------------------------------------------
// Captcha
// -----------------------------------------------------------------------------

var ErrCaptchaNotFound = &AppError{Code: CodeCaptchaNotFound, Kind: KindNotFound, Message: "captcha not found or expired"}
var ErrCaptchaExists = &AppError{Code: CodeCaptchaAlreadyExists, Kind: KindConflict, Message: "captcha already exists"}
var ErrCaptchaGenerationFailed = &AppError{Code: CodeCaptchaGenerationFailed, Kind: KindInternal, Message: "failed to generate captcha"}

// -----------------------------------------------------------------------------
// User Credentials & OTP
// -----------------------------------------------------------------------------

var ErrInvalidPassword = &AppError{Code: CodeInvalidPassword, Kind: KindValidation, Message: "invalid user password"}
var ErrInvalidPhone = &AppError{Code: CodeInvalidPhone, Kind: KindValidation, Message: "invalid user phone number"}
var ErrNotVerifiedPhoneNumber = &AppError{Code: CodeNotVerifiedPhoneNumber, Kind: KindForbidden, Message: "phone number is not verified"}
var ErrInvalidOTP = &AppError{Code: CodeInvalidOTP, Kind: KindValidation, Message: "invalid otp code"}
var ErrOTPExpired = &AppError{Code: CodeOTPExpired, Kind: KindValidation, Message: "otp code has expired"}
var ErrOtpKeyNotFound = &AppError{Code: CodeOtpKeyNotFound, Kind: KindNotFound, Message: "otp key not found in cache"}

// -----------------------------------------------------------------------------
// Users
// -----------------------------------------------------------------------------

var ErrUserNotFound = &AppError{Code: CodeUserNotFound, Kind: KindNotFound, Message: "user not found"}
var ErrUsersNotFound = &AppError{Code: CodeUsersNotFound, Kind: KindNotFound, Message: "no users found matching criteria"}
var ErrUserAlreadyExists = &AppError{Code: CodeUserAlreadyExists, Kind: KindConflict, Message: "user already exists"}

// -----------------------------------------------------------------------------
// RBAC (Roles & Permissions)
// -----------------------------------------------------------------------------

var ErrRoleNotFound = &AppError{Code: CodeRoleNotFound, Kind: KindNotFound, Message: "role not found"}
var ErrRolesNotFound = &AppError{Code: CodeRolesNotFound, Kind: KindNotFound, Message: "no roles found matching criteria"}
var ErrRoleAlreadyExists = &AppError{Code: CodeRoleAlreadyExists, Kind: KindConflict, Message: "role already exists"}
var ErrPermissionNotFound = &AppError{Code: CodePermissionNotFound, Kind: KindNotFound, Message: "permission not found"}
var ErrPermissionsNotFound = &AppError{Code: CodePermissionsNotFound, Kind: KindNotFound, Message: "no permissions found matching criteria"}
var ErrPermissionAlreadyExists = &AppError{Code: CodePermissionAlreadyExists, Kind: KindConflict, Message: "permission already exists"}
var ErrUserPermissionNotFound = &AppError{Code: CodeUserPermissionNotFound, Kind: KindNotFound, Message: "user permission not found in cache"}
var ErrRolePermissionNotFound = &AppError{Code: CodeRolePermissionNotFound, Kind: KindNotFound, Message: "role-permission association not found"}
var ErrUserRoleNotFound = &AppError{Code: CodeUserRoleNotFound, Kind: KindNotFound, Message: "user-role association not found"}
var ErrAssociationFailed = &AppError{Code: CodeAssociationFailed, Kind: KindInternal, Message: "failed to process association operation"}

// -----------------------------------------------------------------------------
// Wallet
// -----------------------------------------------------------------------------

var ErrInvalidAmount = &AppError{Code: CodeInvalidAmount, Kind: KindValidation, Message: "invalid transaction amount"}
var ErrInsufficientBalance = &AppError{Code: CodeInsufficientBalance, Kind: KindValidation, Message: "insufficient wallet balance"}
var ErrReachManBalance = &AppError{Code: CodeReachMaxBalance, Kind: KindValidation, Message: "are you reach or something?!"}

// -----------------------------------------------------------------------------
// Movies & Entities
// -----------------------------------------------------------------------------

var ErrMovieNotFound = &AppError{Code: CodeMovieNotFound, Kind: KindNotFound, Message: "movie not found"}
var ErrMovieAlreadyExists = &AppError{Code: CodeMovieAlreadyExists, Kind: KindConflict, Message: "movie already exists"}
var ErrGenreNotFound = &AppError{Code: CodeGenreNotFound, Kind: KindNotFound, Message: "genre not found"}
var ErrGenreAlreadyExists = &AppError{Code: CodeGenreAlreadyExists, Kind: KindConflict, Message: "genre already exists"}
var ErrPersonNotFound = &AppError{Code: CodePersonNotFound, Kind: KindNotFound, Message: "person not found"}
var ErrPersonAlreadyExists = &AppError{Code: CodePersonAlreadyExists, Kind: KindConflict, Message: "person already exists"}

// -----------------------------------------------------------------------------
// Cinema & Infrastructure
// -----------------------------------------------------------------------------

var ErrCinemaNotFound = &AppError{Code: CodeCinemaNotFound, Kind: KindNotFound, Message: "cinema not found"}
var ErrCinemaAlreadyExists = &AppError{Code: CodeCinemaAlreadyExists, Kind: KindConflict, Message: "cinema already exists"}
var ErrFacilityNotFound = &AppError{Code: CodeFacilityNotFound, Kind: KindNotFound, Message: "facility not found"}
var ErrFacilityAlreadyExists = &AppError{Code: CodeFacilityAlreadyExists, Kind: KindConflict, Message: "facility already exists"}
var ErrCityNotFound = &AppError{Code: CodeCityNotFound, Kind: KindNotFound, Message: "city not found"}
var ErrCityAlreadyExists = &AppError{Code: CodeCityAlreadyExists, Kind: KindConflict, Message: "city already exists"}
var ErrProvinceNotFound = &AppError{Code: CodeProvinceNotFound, Kind: KindNotFound, Message: "province not found"}
var ErrProvinceAlreadyExists = &AppError{Code: CodeProvinceAlreadyExists, Kind: KindConflict, Message: "province already exists"}

// -----------------------------------------------------------------------------
// Comments, Favorites & Ratings
// -----------------------------------------------------------------------------

var ErrCommentNotFound = &AppError{Code: CodeCommentNotFound, Kind: KindNotFound, Message: "comment not found"}
var ErrCommentAlreadyExists = &AppError{Code: CodeCommentAlreadyExists, Kind: KindConflict, Message: "comment already exists"}
var ErrFavoriteNotFound = &AppError{Code: CodeFavoriteNotFound, Kind: KindNotFound, Message: "favorite not found"}
var ErrFavoriteAlreadyExists = &AppError{Code: CodeFavoriteAlreadyExists, Kind: KindConflict, Message: "favorite already exists"}
var ErrRatingNotFound = &AppError{Code: CodeRatingNotFound, Kind: KindNotFound, Message: "rating not found"}
var ErrRatingAlreadyExists = &AppError{Code: CodeRatingAlreadyExists, Kind: KindConflict, Message: "rating already exists"}

// -----------------------------------------------------------------------------
// System, Media & Database
// -----------------------------------------------------------------------------

var ErrDuplicateKey = &AppError{Code: CodeDuplicateKey, Kind: KindConflict, Message: "record already exists"}
var ErrResourceInUse = &AppError{Code: CodeResourceInUse, Kind: KindConflict, Message: "resource is in use and cannot be modified or deleted"}
var ErrInvalidFileType = &AppError{Code: CodeInvalidFileType, Kind: KindValidation, Message: "invalid file type provided"}
var ErrFileSizeTooLarge = &AppError{Code: CodeFileSizeTooLarge, Kind: KindValidation, Message: "file size exceeds maximum allowed limit"}
var ErrDatabase = &AppError{Code: CodeDatabase, Kind: KindDatabase, Message: "database operation failed"}
var ErrCache = &AppError{Code: CodeCache, Kind: KindInternal, Message: "cache operation failed"}
var ErrInternal = &AppError{Code: CodeInternal, Kind: KindInternal, Message: "internal server error"}

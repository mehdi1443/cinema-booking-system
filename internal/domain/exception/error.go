package exception

import (
	"errors"
	"fmt"
)

const (
	ErrorTypeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrorTypeExpiredToken       = "EXPIRED_TOKEN"
	ErrorTypeInvalidToken       = "INVALID_TOKEN"
	ErrorTypeUnauthorized       = "UNAUTHORIZED"
)

type AuthError struct {
	Type        string
	Message     string
	OriginalErr error
}

func (e AuthError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.OriginalErr.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e AuthError) Unwrap() error {
	return e.OriginalErr
}

func NewInvalidCredentialsError(message string, originalErr error) *AuthError {
	if message == "" {
		message = "username and password not match"
	}
	return &AuthError{
		Type:        ErrorTypeInvalidCredentials,
		Message:     message,
		OriginalErr: originalErr,
	}
}

func NewExpiredTokenError(originalErr error) *AuthError {
	return &AuthError{
		Type:        ErrorTypeExpiredToken,
		Message:     "Authentication token has expired",
		OriginalErr: originalErr,
	}
}

func NewInvalidTokenError(originalErr error) *AuthError {
	return &AuthError{
		Type:        ErrorTypeInvalidToken,
		Message:     "Invalid authentication token",
		OriginalErr: originalErr,
	}
}

func NewUnauthorizedError(message string, originalErr error) *AuthError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &AuthError{
		Type:        ErrorTypeUnauthorized,
		Message:     message,
		OriginalErr: originalErr,
	}
}

func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

func GetAuthErrorType(err error) string {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type
	}
	return ""
}

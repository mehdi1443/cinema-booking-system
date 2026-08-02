package service_contract

import "errors"

//TODO : error handling in services
var (
	// jwt service errors
	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid or malformed token")
)

const (
	UserIDKeyMapClaims    string = "user_id"
	TokenTypeKeyMapClaims string = "token_type"
	IssuedAtKeyMapClaims  string = "iat"
	ExpireAtKeyMapClaims  string = "exp"
)

package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

type contexttokenService struct {
}

func NewContextTokenService() service_contract.ContextTokenServiceInterface {
	return &contexttokenService{}
}

var (
	ErrUserIDNotFoundInContext      = errors.New("user ID not found in context")
	ErrInvalidUserIDType            = errors.New("user ID in context has invalid type")
	ErrPermissionsNotFoundInContext = errors.New("user permissions not found in context")
	ErrInvalidPermissionsType       = errors.New("user permissions in context have invalid type")

	ErrAccessTokenNotFoundInContext  = errors.New("access token not found in context")
	ErrInvalidAccessTokenType        = errors.New("access token in context has invalid type")
	ErrInvalidRefreshTokenType       = errors.New("refresh token in context has invalid type")
	ErrRefreshTokenNotFoundInContext = errors.New("refresh token not found in context")

	ErrAccessTokenExpNotFoundInContext = errors.New("access token expiration time not found in context")
	ErrInvalidAccessTokenExpType       = errors.New("access token expiration time in context has invalid type")

	ErrRefreshTokenExpNotFoundInContext = errors.New("refresh token expiration time not found in context")
	ErrInvalidRefreshTokenExpType       = errors.New("refresh token expiration time in context has invalid type")
)

func (s *contexttokenService) GetTokenFromGinContext(gctx *gin.Context) (dto.TokenInContext, error) {
	if gctx == nil {
		return dto.TokenInContext{}, errors.New("gin context is nil")
	}

	// 1. Extract User ID
	rawUserID := getContextValue(gctx, domainjwt.CtxUserIDKey, "user_id")
	if rawUserID == nil {
		return dto.TokenInContext{}, ErrUserIDNotFoundInContext
	}

	userID, ok := rawUserID.(int64)
	if !ok {
		return dto.TokenInContext{}, ErrInvalidUserIDType
	}

	// 2. Extract Permissions
	rawPermissions := getContextValue(gctx, domainjwt.CtxPermissionsKey, "permissions")
	if rawPermissions == nil {
		return dto.TokenInContext{}, ErrPermissionsNotFoundInContext
	}

	permissionNames, ok := rawPermissions.([]string)
	if !ok {
		return dto.TokenInContext{}, ErrInvalidPermissionsType
	}

	// 3. Construct and return DTO
	return dto.TokenInContext{
		UserID:          userID,
		PermissionNames: permissionNames,
	}, nil
}

func (s *contexttokenService) GetJWTFromGinContext(gctx *gin.Context) (dto.JWTTokenInContext, error) {
	if gctx == nil {
		return dto.JWTTokenInContext{}, errors.New("gin context is nil")
	}

	// 1. Extract User ID
	rawT := getContextValue(gctx, domainjwt.CtxTheTokenKey)
	if rawT == nil {
		return dto.JWTTokenInContext{}, ErrAccessTokenNotFoundInContext
	}
	T, ok := rawT.(string)
	if !ok {
		return dto.JWTTokenInContext{}, ErrInvalidAccessTokenType
	}

	rawTXTime := getContextValue(gctx, domainjwt.CtxTheTokenExpireAtKey)
	if rawTXTime == nil {
		return dto.JWTTokenInContext{}, ErrAccessTokenExpNotFoundInContext
	}
	TXTime, ok := rawTXTime.(time.Time)
	if !ok {
		return dto.JWTTokenInContext{}, ErrInvalidAccessTokenExpType
	}

	// 2. Construct and return DTO
	return dto.JWTTokenInContext{
		TheToken:         T,
		TheTokenExpireAt: TXTime,
	}, nil
}

// GetUserIDFromGinContext extracts only the User ID from context
func (s *contexttokenService) GetUserIDFromGinContext(gctx *gin.Context) (int64, error) {
	if gctx == nil {
		return 0, errors.New("gin context is nil")
	}

	rawUserID := getContextValue(gctx, domainjwt.CtxUserIDKey, "user_id")
	if rawUserID == nil {
		return 0, ErrUserIDNotFoundInContext
	}

	userID, ok := rawUserID.(int64)
	if !ok {
		return 0, ErrInvalidUserIDType
	}

	return userID, nil
}

// GetPermissionFromGinContext extracts only the Permissions slice from context
func (s *contexttokenService) GetPermissionFromGinContext(gctx *gin.Context) ([]string, error) {
	if gctx == nil {
		return nil, errors.New("gin context is nil")
	}

	rawPermissions := getContextValue(gctx, domainjwt.CtxPermissionsKey, "permissions")
	if rawPermissions == nil {
		return nil, ErrPermissionsNotFoundInContext
	}

	permissionNames, ok := rawPermissions.([]string)
	if !ok {
		return nil, ErrInvalidPermissionsType
	}

	return permissionNames, nil
}

// Helper to look up values across Gin's c.Set() map AND Request Context
func getContextValue(gctx *gin.Context, keys ...any) any {
	for _, key := range keys {
		if strKey, ok := key.(string); ok {
			if val, exists := gctx.Get(strKey); exists && val != nil {
				return val
			}
		}

		if val := gctx.Value(key); val != nil {
			return val
		}

		if gctx.Request != nil && gctx.Request.Context() != nil {
			if val := gctx.Request.Context().Value(key); val != nil {
				return val
			}
		}
	}
	return nil
}

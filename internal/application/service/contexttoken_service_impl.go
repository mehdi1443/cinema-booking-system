package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
)

type contexttokenService struct {
}

func NewContextTokenService() service_contract.ContextTokenServiceInterface {
	return &contexttokenService{}
}

var (
	ErrUserIDNotFoundInContext = errors.New("user ID not found in context")
	ErrInvalidUserIDType       = errors.New("user ID in context has invalid type")
	ErrRolesNotFoundInContext  = errors.New("user roles not found in context")
	ErrInvalidRolesType        = errors.New("user roles in context have invalid type")
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

	// 2. Extract Roles
	rawRoles := getContextValue(gctx, domainjwt.CtxRolesKey, "roles")
	if rawRoles == nil {
		return dto.TokenInContext{}, ErrRolesNotFoundInContext
	}

	roles, ok := rawRoles.([]entity.Role)
	if !ok {
		return dto.TokenInContext{}, ErrInvalidRolesType
	}

	// 3. Construct and return DTO
	return dto.TokenInContext{
		UserID: userID,
		Roles:  roles,
	}, nil
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

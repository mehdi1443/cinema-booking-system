package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func loadUserPermissions(
	c *gin.Context,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contextTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	appLogger domainlogger.Logger,
) (int64, []string, bool) {
	ctx := c.Request.Context()
	userID, err := contextTokenService.GetUserIDFromGinContext(c)
	if err != nil {
		response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeTokenContext), exception.ErrGetTokenFromContext.Error())
		c.Abort()
		return 0, nil, false
	}
	user, err := userRepo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeUnauthorized), "User account not found")
		} else {
			appLogger.Error(ctx, "failed to load authenticated user", domainlogger.Int64("user_id", userID), domainlogger.Error(err))
			response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "Failed to verify user account")
		}
		c.Abort()
		return 0, nil, false
	}
	if user == nil {
		response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeUnauthorized), "User account not found")
		c.Abort()
		return 0, nil, false
	}
	if user.Status != enum.UserActive {
		response.ErrorWithCode(c, http.StatusForbidden, string(exception.CodeForbidden), "User account is not active")
		c.Abort()
		return 0, nil, false
	}
	if !user.IsPhoneVerified {
		response.ErrorWithCode(c, http.StatusUnauthorized, string(exception.CodeNotVerifiedPhoneNumber), "Phone number is not verified")
		c.Abort()
		return 0, nil, false
	}

	userIDString := strconv.FormatInt(userID, 10)
	permissions, err := userPermissionCacheRepo.GetUserPermissions(ctx, userIDString)
	if errors.Is(err, exception.ErrUserPermissionNotFound) {
		permissions, err = rbacRepo.GetUserPermissionNames(ctx, userID)
		if err == nil {
			if cacheErr := userPermissionCacheRepo.SaveUserPermissions(ctx, userIDString, permissions, userCacheConfig.ExpirationDuration()); cacheErr != nil {
				appLogger.Warn(ctx, "failed to cache permissions", domainlogger.Error(cacheErr))
			}
		}
	}
	if err != nil {
		appLogger.Error(ctx, "failed to load user permissions", domainlogger.Error(err))
		response.ErrorWithCode(c, http.StatusInternalServerError, string(exception.CodeInternal), "Failed to verify access permissions")
		c.Abort()
		return 0, nil, false
	}
	c.Set(domainjwt.CtxPermissionsKey, permissions)
	return userID, permissions, true
}

func hasPermission(permissions []string, required string) bool {
	for _, permission := range permissions {
		if permission == required || permission == enum.PermissionAll.String() {
			return true
		}
	}
	return false
}

func ControlUserAccess(
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contextTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	requiredPermission string,
	appLogger domainlogger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, permissions, ok := loadUserPermissions(c, rbacRepo, userPermissionCacheRepo, contextTokenService, userRepo, userCacheConfig, appLogger)
		if !ok {
			return
		}
		if !hasPermission(permissions, requiredPermission) {
			response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
			c.Abort()
			return
		}
		c.Next()
	}
}

func ControlSelfOrPermission(
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contextTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	selfPermission string,
	anyUserPermission string,
	paramName string,
	appLogger domainlogger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, permissions, ok := loadUserPermissions(c, rbacRepo, userPermissionCacheRepo, contextTokenService, userRepo, userCacheConfig, appLogger)
		if !ok {
			return
		}
		targetID, err := strconv.ParseInt(c.Param(paramName), 10, 64)
		if err != nil || targetID <= 0 {
			response.Error(c, http.StatusBadRequest, "Invalid user id")
			c.Abort()
			return
		}
		required := anyUserPermission
		if targetID == userID {
			required = selfPermission
		}
		if !hasPermission(permissions, required) {
			response.Error(c, http.StatusForbidden, "You do not have permission to access this user")
			c.Abort()
			return
		}
		c.Next()
	}
}

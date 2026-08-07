package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func ControlUserAccess(
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	requiredPermission string,
	appLogger logger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := contexTokenService.GetUserIDFromGinContext(c)
		if err != nil {
			appLogger.Error(logger.Auth, logger.RBAC, "failed to extract token context", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
			c.Abort()
			return
		}

		user, err := userRepo.FindUserByID(c.Request.Context(), userID)
		if err != nil || user == nil {
			appLogger.Warn(logger.Auth, logger.RBAC, "user account not found", map[logger.ExtraKey]interface{}{
				logger.UserID: userID,
			})
			response.Error(c, http.StatusUnauthorized, "User account not found")
			c.Abort()
			return
		}

		if !user.IsPhoneVerified {
			response.Error(c, http.StatusUnauthorized, "Phone number is not verified")
			c.Abort()
			return
		}

		userIDStr := strconv.Itoa(int(userID))
		permissions, err := userPermissionCacheRepo.GetUserPermissions(c.Request.Context(), userIDStr)
		if err != nil {
			if errors.Is(err, exception.ErrUserPermissionNotFound) {
				permissions, err = rbacRepo.GetUserPermissionNames(c.Request.Context(), userID)
				if err != nil {
					appLogger.Error(logger.Auth, logger.RBAC, "failed to fetch permissions from DB", map[logger.ExtraKey]interface{}{
						logger.UserID:       userID,
						logger.ErrorMessage: err.Error(),
					})
					response.Error(c, http.StatusInternalServerError, "Failed to verify access permissions")
					c.Abort()
					return
				}

				if cacheErr := userPermissionCacheRepo.SaveUserPermissions(c.Request.Context(), userIDStr, permissions, userCacheConfig.ExpirationDuration()); cacheErr != nil {
					appLogger.Warn(logger.Auth, logger.RBAC, "failed to save permissions to cache", map[logger.ExtraKey]interface{}{
						logger.UserID:       userID,
						logger.ErrorMessage: cacheErr.Error(),
					})
				}
			} else {
				appLogger.Error(logger.Auth, logger.RBAC, "error reading permission cache", map[logger.ExtraKey]interface{}{
					logger.UserID:       userID,
					logger.ErrorMessage: err.Error(),
				})
				response.Error(c, http.StatusInternalServerError, "Internal authorization error")
				c.Abort()
				return
			}
		}

		hasPermission := false
		for _, perm := range permissions { //FIXME
			if perm == requiredPermission || perm == "ALL" {
				hasPermission = true
				break
			}
		}

		if !hasPermission { //FIXME
			appLogger.Warn(logger.Auth, logger.RBAC, "access denied due to missing permission", map[logger.ExtraKey]interface{}{
				logger.UserID: userID,
			})
			response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set(domainjwt.CtxPermissionsKey, permissions)
		c.Next()
	}
}

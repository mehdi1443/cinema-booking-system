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
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

func ControlUserAccess(
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	requiredPermission string,
	appLogger domainlogger.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID, err := contexTokenService.GetUserIDFromGinContext(c)
		if err != nil {
			appLogger.Error(ctx, "failed to extract token context", domainlogger.Error(err))
			response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
			c.Abort()
			return
		}

		user, err := userRepo.FindUserByID(ctx, userID)
		if err != nil || user == nil {
			appLogger.Warn(ctx, "user account not found", domainlogger.Int64("user_id", userID))
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
		permissions, err := userPermissionCacheRepo.GetUserPermissions(ctx, userIDStr)
		if err != nil {
			if errors.Is(err, exception.ErrUserPermissionNotFound) {
				permissions, err = rbacRepo.GetUserPermissionNames(ctx, userID)
				if err != nil {
					appLogger.Error(ctx, "failed to fetch permissions from DB",
						domainlogger.Int64("user_id", userID),
						domainlogger.Error(err),
					)
					response.Error(c, http.StatusInternalServerError, "Failed to verify access permissions")
					c.Abort()
					return
				}

				if cacheErr := userPermissionCacheRepo.SaveUserPermissions(ctx, userIDStr, permissions, userCacheConfig.ExpirationDuration()); cacheErr != nil {
					appLogger.Warn(ctx, "failed to save permissions to cache",
						domainlogger.Int64("user_id", userID),
						domainlogger.Error(cacheErr),
					)
				}
			} else {
				appLogger.Error(ctx, "error reading permission cache",
					domainlogger.Int64("user_id", userID),
					domainlogger.Error(err),
				)
				response.Error(c, http.StatusInternalServerError, "Internal authorization error")
				c.Abort()
				return
			}
		}

		hasPermission := false
		for _, perm := range permissions {
			if perm == requiredPermission || perm == "ALL" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			appLogger.Warn(ctx, "access denied due to missing permission", domainlogger.Int64("user_id", userID))
			response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set(domainjwt.CtxPermissionsKey, permissions)
		c.Next()
	}
}

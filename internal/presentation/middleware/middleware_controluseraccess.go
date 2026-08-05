package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainjwt "github.com/mehdi1443/cinema-booking-system/internal/domain/jwt"
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
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := contexTokenService.GetUserIDFromGinContext(c)
		if err != nil {
			log.Printf("[AUTH ERROR] Failed to extract token context: %v", err)
			response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error())
			c.Abort()
			return
		}

		user, err := userRepo.FindUserByID(c.Request.Context(), userID)
		if err != nil || user == nil {
			log.Printf("[AUTH ERROR] User not found: %v", err)
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
		// log.Println(permissions, "**",err)
		if err != nil {
			if errors.Is(err, exception.ErrUserPermissionNotFound) {
				permissions, err = rbacRepo.GetUserPermissionNames(c.Request.Context(), userID)
				// log.Println(permissions, "*")
				if err != nil {
					log.Printf("[AUTH ERROR] Failed to fetch permissions from DB: %v", err)
					response.Error(c, http.StatusInternalServerError, "Failed to verify access permissions")
					c.Abort()
					return
				}

				if cacheErr := userPermissionCacheRepo.SaveUserPermissions(c.Request.Context(), userIDStr, permissions, userCacheConfig.ExpirationDuration()); cacheErr != nil {
					log.Printf("[CACHE WARN] Failed to save permissions to cache: %v", cacheErr)
				}
			} else {
				log.Printf("[CACHE ERROR] Error reading permission cache: %v", err)
				response.Error(c, http.StatusInternalServerError, "Internal authorization error")
				c.Abort()
				return
			}
		}

		hasPermission := false
		for _, perm := range permissions { // FIXME!!!
			if perm == requiredPermission || perm == "ALL" {
				hasPermission = true
			}
		}
		hasPermission = true //FIXME
		if !hasPermission {
			response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set(domainjwt.CtxPermissionsKey, permissions)
		c.Next()
	}
}

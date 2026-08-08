package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

func SetUpRBACRouter(
	router *gin.RouterGroup,
	rbacHandler *handler.RBACHandler,
	authMiddleware gin.HandlerFunc,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	appLogger logger.Logger,
) {

	// Role Management
	roles := router.Group("/roles")
	roles.Use(limiter.PerMinute(limiterConfig.HighRate))
	roles.Use(authMiddleware)
	{
		roles.GET("/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "roles:read", appLogger),
			rbacHandler.GetRoleByID,
		)
		roles.POST("/",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "roles:read", appLogger),
			rbacHandler.CreateRole,
		)
		roles.PUT("/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "roles:read", appLogger),
			rbacHandler.UpdateRole,
		)
		roles.POST("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "roles:append_permissions", appLogger),
			rbacHandler.AppendPermissionsToRole,
		)
		roles.DELETE("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "roles:revoke_permissions", appLogger),
			rbacHandler.DeletePermissionsOfRole,
		)
	}

	// Permission Management
	permissions := router.Group("/permissions")
	permissions.Use(limiter.PerMinute(limiterConfig.HighRate))
	permissions.Use(authMiddleware)
	{
		permissions.GET("/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "permissions:read", appLogger),
			rbacHandler.GetPermissionByID,
		)
		permissions.GET("/role/:name",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "permissions:read_by_role", appLogger),
			rbacHandler.GetPermissionsByRoleName,
		)
	}
}

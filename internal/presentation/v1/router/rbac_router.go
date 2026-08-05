package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRBACRouter(
	router *gin.RouterGroup,
	rbacHandler *handler.RBACHandler,
	authMiddleware gin.HandlerFunc,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	rbacService service_contract.RBACServiceInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userService service_contract.UserServiceInterface,
) {

	// Role Management
	roles := router.Group("/roles")
	roles.Use(limiter.PerMinute(limiterConfig.HighRate), authMiddleware)
	{
		roles.GET("/:id",
			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "roles:read"),
			rbacHandler.GetRoleByID,
		)
		roles.POST("/:id/permissions",
			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "roles:assign_permissions"),
			rbacHandler.AssignPermissionsToRole,
		)
		roles.DELETE("/:id/permissions",
			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "roles:revoke_permissions"),
			rbacHandler.DischargePermissionsOfRole,
		)
	}

	// Permission Management
	permissions := router.Group("/permissions")
	permissions.Use(limiter.PerMinute(1), authMiddleware)
	{
		permissions.GET("/:id",
			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "permissions:read"),
			rbacHandler.GetPermissionByID,
		)
		permissions.GET("/role/:name",
			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "permissions:read_by_role"),
			rbacHandler.GetPermissionsByRoleName,
		)
	}
}

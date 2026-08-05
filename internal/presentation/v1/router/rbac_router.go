package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRBACRouter(
	router *gin.RouterGroup,
	rbacHandler *handler.RBACHandler,
	authMiddleware gin.HandlerFunc,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	rbacRepo DB.RBACInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
) {

	// Role Management
	roles := router.Group("/roles")
	roles.Use(limiter.PerMinute(limiterConfig.HighRate), authMiddleware)
	{
		roles.GET("/:id",
			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "roles:read"),
			rbacHandler.GetRoleByID,
		)
		roles.POST("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "roles:assign_permissions"),
			rbacHandler.AssignPermissionsToRole,
		)
		roles.DELETE("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "roles:revoke_permissions"),
			rbacHandler.DischargePermissionsOfRole,
		)
	}

	// Permission Management
	permissions := router.Group("/permissions")
	permissions.Use(limiter.PerMinute(1), authMiddleware)
	{
		permissions.GET("/:id",
			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "permissions:read"),
			rbacHandler.GetPermissionByID,
		)
		permissions.GET("/role/:name",
			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "permissions:read_by_role"),
			rbacHandler.GetPermissionsByRoleName,
		)
	}
}

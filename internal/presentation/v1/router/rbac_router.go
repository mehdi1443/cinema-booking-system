package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRBACRouter(router *gin.RouterGroup, rbacHandler *handler.RBACHandler,
	authMiddleware gin.HandlerFunc) {
	// // Role Management
	roles := router.Group("/roles")
	{
		roles.GET("/:id", rbacHandler.GetRoleByID)
		roles.POST("/:id/permissions", rbacHandler.AssignPermissionsToRole)
		roles.DELETE("/:id/permissions", rbacHandler.DischargePermissionsOfRole)
	}

	// Permission Management
	permissions := router.Group("/permissions")
	{
		permissions.GET("/:id", rbacHandler.GetPermissionByID)
		permissions.GET("/role/:name", rbacHandler.GetPermissionsByRoleName)
	}
}

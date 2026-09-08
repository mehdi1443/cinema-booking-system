package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRBACRouter(
	router *gin.RouterGroup,
	rbacHandler *handler.RBACHandler,
) {
	roles := router.Group("/roles")
	{
		roles.GET("/:id", rbacHandler.GetRoleByID)
		roles.POST("/", rbacHandler.CreateRole)
		roles.PUT("/:id", rbacHandler.UpdateRole)
		roles.POST("/:id/permissions", rbacHandler.AppendPermissionsToRole)
		roles.DELETE("/:id/permissions", rbacHandler.DeletePermissionsOfRole)
	}

	permissions := router.Group("/permissions")
	{
		permissions.GET("/:id", rbacHandler.GetPermissionByID)
		permissions.GET("/role/:name", rbacHandler.GetPermissionsByRoleName)
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRBACRouter(
	router *gin.RouterGroup,
	rbacHandler *handler.RBACHandler,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
) {
	roles := router.Group("/roles")
	roles.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		roles.GET("/:id", rbacHandler.GetRoleByID)
		roles.POST("/", rbacHandler.CreateRole)
		roles.PUT("/:id", rbacHandler.UpdateRole)
		roles.POST("/:id/permissions", rbacHandler.AppendPermissionsToRole)
		roles.DELETE("/:id/permissions", rbacHandler.DeletePermissionsOfRole)
	}

	permissions := router.Group("/permissions")
	permissions.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		permissions.GET("/:id", rbacHandler.GetPermissionByID)
		permissions.GET("/role/:name", rbacHandler.GetPermissionsByRoleName)
	}
}

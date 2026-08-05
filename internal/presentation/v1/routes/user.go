package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

// TODO : ignrore fileds for base response in each endpoints + permission&Accessibility + userpermission cache issue
func RegisterRoutes(
	rg *gin.RouterGroup,
	userHandler *handler.UserHandler,
	rbacHandler *handler.RBACHandler,
	walletHandler *handler.WalletHandler,
) {
	// User Endpoints
	users := rg.Group("/users")
	{
		//* specific
		users.PATCH("/:id/password", userHandler.UpdatePassword)	

		//profile
		users.GET("/:id/profile", userHandler.GetProfile)
		users.PUT("/:id/profile", userHandler.UpdateProfile)

		//status
		users.PATCH("/:id/status", userHandler.ChangeStatus)

		//roles
		users.GET("/:id/roles", rbacHandler.GetUserRoles)
		users.PUT("/:id/roles", userHandler.UpdateUserRoles)

		//permission
		users.GET("/:id/permissions", rbacHandler.GetUserPermissions)
		users.DELETE("/:id", userHandler.DeleteUserByID)

		//* all
		users.GET("", userHandler.GetUsers) // query

	}

	// // Role Management
	roles := rg.Group("/roles")
	{
		roles.GET("/:id", rbacHandler.GetRoleByID)
		roles.POST("/:id/permissions", rbacHandler.AssignPermissionsToRole)
		roles.DELETE("/:id/permissions", rbacHandler.DischargePermissionsOfRole)
	}

	// Permission Management
	permissions := rg.Group("/permissions")
	{
		permissions.GET("/:id", rbacHandler.GetPermissionByID)
		permissions.GET("/role/:name", rbacHandler.GetPermissionsByRoleName)
	}

	wallet := rg.Group("/users/:id/wallet")
	{
		wallet.GET("/balance", walletHandler.GetBalance)
		wallet.POST("/balance", walletHandler.ChangeBalance)
	}
}

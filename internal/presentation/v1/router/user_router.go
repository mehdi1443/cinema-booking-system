package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

// TODO : ignrore fileds for base response in each endpoints + permission&Accessibility + userpermission cache issue
func SetUpUserRouter(
	rg *gin.RouterGroup,
	userHandler *handler.UserHandler,
	rbacHandler *handler.RBACHandler,
	walletHandler *handler.WalletHandler,
	authMiddleware  gin.HandlerFunc,
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

		//* wallet
		users.GET("/:id/wallet/balance", walletHandler.GetBalance)
		users.POST("/:id/wallet/balance", walletHandler.ChangeBalance)
	}

}

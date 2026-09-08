package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpUserRouter(
	rg *gin.RouterGroup,
	userHandler *handler.UserHandler,
	walletHandler *handler.WalletHandler,
) {
	users := rg.Group("/users")
	{
		users.PATCH("/:id/password", userHandler.UpdatePassword)
		users.GET("/:id/profile", userHandler.GetProfile)
		users.PUT("/:id/profile", userHandler.UpdateProfile)
		users.PATCH("/:id/status", userHandler.ChangeStatus)
		users.GET("/:id/roles", userHandler.GetUserRoles)
		users.PUT("/:id/roles", userHandler.UpdateUserRoles)
		users.GET("/:id/permissions", userHandler.GetUserPermissions)
		users.DELETE("/:id", userHandler.DeleteUserByID)
		users.GET("", userHandler.GetUsers)
		users.GET("/:id/wallet/balance", walletHandler.GetBalance)
		users.POST("/:id/wallet/balance", walletHandler.ChangeBalance)
	}
}

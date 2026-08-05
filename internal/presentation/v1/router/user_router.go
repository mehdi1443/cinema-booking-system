package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

// TODO : ignrore fileds for base response in each endpoints + permission&Accessibility + userpermission cache issue
func SetUpUserRouter(
	rg *gin.RouterGroup,
	userHandler *handler.UserHandler,
	rbacHandler *handler.RBACHandler,
	walletHandler *handler.WalletHandler,
	authMiddleware gin.HandlerFunc,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	rbacService service_contract.RBACServiceInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userService service_contract.UserServiceInterface,
) {
	// User Endpoints
	users := rg.Group("/users")
	users.Use(limiter.PerMinute(limiterConfig.HighRate), authMiddleware)
	{
		//* specific
		users.PATCH("/:id/password",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:update_password"),
			userHandler.UpdatePassword,
		)

		//profile
		users.GET("/:id/profile",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:read_profile"),
			userHandler.GetProfile,
		)
		users.PUT("/:id/profile",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:update_profile"),
			userHandler.UpdateProfile,
		)

		//status
		users.PATCH("/:id/status",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:change_status"),
			userHandler.ChangeStatus,
		)

		//roles
		users.GET("/:id/roles",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:read_roles"),
			userHandler.GetUserRoles,
		)
		users.PUT("/:id/roles",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:update_roles"),
			userHandler.UpdateUserRoles,
		)

		//permission
		users.GET("/:id/permissions",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:read_permissions"),
			userHandler.GetUserPermissions,
		)
		users.DELETE("/:id",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:delete"),
			userHandler.DeleteUserByID,
		)

		//* all
		users.GET("",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "users:list"),
			userHandler.GetUsers,
		) // query

		//* wallet
		users.GET("/:id/wallet/balance",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "wallet:read_balance"),
			walletHandler.GetBalance,
		)
		users.POST("/:id/wallet/balance",

			middleware.ControlUserAccess(rbacService, contexTokenService, userService, "wallet:change_balance"),
			walletHandler.ChangeBalance,
		)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
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
	rbacRepo DB.RBACInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
) {
	// User Endpoints
	users := rg.Group("/users")
	users.Use(limiter.PerMinute(limiterConfig.HighRate), authMiddleware)
	{
		//* specific
		users.PATCH("/:id/password",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:update_password"),
			userHandler.UpdatePassword,
		)

		//profile
		users.GET("/:id/profile",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:read_profile"),
			userHandler.GetProfile,
		)
		users.PUT("/:id/profile",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:update_profile"),
			userHandler.UpdateProfile,
		)

		//status
		users.PATCH("/:id/status",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:change_status"),
			userHandler.ChangeStatus,
		)

		//roles
		users.GET("/:id/roles",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:read_roles"),
			userHandler.GetUserRoles,
		)
		users.PUT("/:id/roles",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:update_roles"),
			userHandler.UpdateUserRoles,
		)

		//permission
		users.GET("/:id/permissions",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:read_permissions"),
			userHandler.GetUserPermissions,
		)
		users.DELETE("/:id",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:delete"),
			userHandler.DeleteUserByID,
		)

		//* all
		users.GET("",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "users:list"),
			userHandler.GetUsers,
		) // query

		//* wallet
		users.GET("/:id/wallet/balance",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "wallet:read_balance"),
			walletHandler.GetBalance,
		)
		users.POST("/:id/wallet/balance",

			middleware.ControlUserAccess(rbacRepo, contexTokenService, userRepo, "wallet:change_balance"),
			walletHandler.ChangeBalance,
		)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
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
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
) {
	// User Endpoints
	users := rg.Group("/users")
	users.Use(limiter.PerMinute(limiterConfig.HighRate), authMiddleware)
	{
		//* specific
		users.PATCH("/:id/password",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_password"),
			userHandler.UpdatePassword,
		)

		// profile
		users.GET("/:id/profile",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_profile"),
			userHandler.GetProfile,
		)
		users.PUT("/:id/profile",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_profile"),
			userHandler.UpdateProfile,
		)

		// status
		users.PATCH("/:id/status",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:change_status"),
			userHandler.ChangeStatus,
		)

		// roles
		users.GET("/:id/roles",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_roles"),
			userHandler.GetUserRoles,
		)
		users.PUT("/:id/roles",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_roles"),
			userHandler.UpdateUserRoles,
		)

		// permission
		users.GET("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_permissions"),
			userHandler.GetUserPermissions,
		)
		users.DELETE("/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:delete"),
			userHandler.DeleteUserByID,
		)

		//* all
		users.GET("",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:list"),
			userHandler.GetUsers,
		) // query

		//* wallet
		users.GET("/:id/wallet/balance",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "wallet:read_balance"),
			walletHandler.GetBalance,
		)
		users.POST("/:id/wallet/balance",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "wallet:change_balance"),
			walletHandler.ChangeBalance,
		)
	}
}

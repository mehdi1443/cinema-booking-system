package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
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
	appLogger logger.Logger,
) {
	// User Endpoints
	users := rg.Group("/users")
	users.Use(limiter.PerMinute(limiterConfig.HighRate))
	users.Use(authMiddleware)
	{
		//* specific
		users.PATCH("/:id/password",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_password", appLogger),
			userHandler.UpdatePassword,
		)

		// profile
		users.GET("/:id/profile",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_profile", appLogger),
			userHandler.GetProfile,
		)
		users.PUT("/:id/profile",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_profile", appLogger),
			userHandler.UpdateProfile,
		)

		// status
		users.PATCH("/:id/status",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:change_status", appLogger),
			userHandler.ChangeStatus,
		)

		// roles
		users.GET("/:id/roles",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_roles", appLogger),
			userHandler.GetUserRoles,
		)
		users.PUT("/:id/roles",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:update_roles", appLogger),
			userHandler.UpdateUserRoles,
		)

		// permission
		users.GET("/:id/permissions",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:read_permissions", appLogger),
			userHandler.GetUserPermissions,
		)
		users.DELETE("/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:delete", appLogger),
			userHandler.DeleteUserByID,
		)

		//* all
		users.GET("",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "users:list", appLogger),
			userHandler.GetUsers,
		) // query

		//* wallet
		users.GET("/:id/wallet/balance",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "wallet:read_balance", appLogger),
			walletHandler.GetBalance,
		)
		users.POST("/:id/wallet/balance",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "wallet:change_balance", appLogger),
			walletHandler.ChangeBalance,
		)
	}
}

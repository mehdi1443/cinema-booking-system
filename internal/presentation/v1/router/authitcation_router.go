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

func SetUpAuthRouter(
	router *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
	limiter *middleware.RedisLimiter,
	limiterConfig *bootstrap.Limiter,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
) {
	auth := router.Group("/auth")
	auth.Use(limiter.PerMinute(limiterConfig.HighRate))
	{
		auth.POST("/otp/send", authHandler.SendOTP) // TODO: Captcha Handling
		auth.POST("/otp/verify", authHandler.VerifyOTP)
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.Refresh)

		auth.POST("/logout",
			authMiddleware,
			//! NOTICE: DONT EVEN THINK TO TOUCH IT !!! : middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "auth:logout"),
			authHandler.Logout,
		)

		auth.POST("/verify-phone",
			authMiddleware,
			//! NOTICE: DONT EVEN THINK TO TOUCH IT !!! : middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "auth:verify_phone"),
			authHandler.VerifyPhone,
		)
		auth.POST("/verify-email",
			authMiddleware,
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "auth:verify_email"),
			authHandler.VerifyEmail,
		)
	}
}

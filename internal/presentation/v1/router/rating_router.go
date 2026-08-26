package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRatingRouter(
	rg *gin.RouterGroup,
	ratingHandler *handler.RatingHandler,
	authMiddleware gin.HandlerFunc,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	appLogger domainlogger.Logger,
) {
	// Public Routes
	rg.GET("/movies/:id/ratings", ratingHandler.GetMovieRatings)
	rg.GET("/cinemas/:id/ratings", ratingHandler.GetCinemaRatings)

	// Protected Routes
	protected := rg.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("/ratings",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "ratings:submit", appLogger),
			ratingHandler.SubmitRating,
		)

		protected.GET("/movies/:id/user-rating",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "ratings:read_self", appLogger),
			ratingHandler.GetUserMovieRating,
		)

		protected.GET("/cinemas/:id/user-rating",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "ratings:read_self", appLogger),
			ratingHandler.GetUserCinemaRating,
		)

		protected.DELETE("/ratings/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "ratings:delete", appLogger),
			ratingHandler.DeleteRating,
		)
	}
}

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

func SetUpFavoriteRouter(
	rg *gin.RouterGroup,
	favoriteHandler *handler.FavoriteHandler,
	authMiddleware gin.HandlerFunc,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	appLogger logger.Logger,
) {
	favorites := rg.Group("/favorites")
	favorites.Use(authMiddleware)
	{
		favorites.POST("/movies",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "favorites:add_movie", appLogger),
			favoriteHandler.AddMovieToFavorites,
		)

		favorites.DELETE("/movies/:movie_id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "favorites:remove_movie", appLogger),
			favoriteHandler.RemoveMovieFromFavorites,
		)

		favorites.POST("/cinemas",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "favorites:add_cinema", appLogger),
			favoriteHandler.AddCinemaToFavorites,
		)

		favorites.DELETE("/cinemas/:cinema_id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "favorites:remove_cinema", appLogger),
			favoriteHandler.RemoveCinemaFromFavorites,
		)

		favorites.GET("",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "favorites:read", appLogger),
			favoriteHandler.GetUserFavorites,
		)
	}
}

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpFavoriteRouter(rg *gin.RouterGroup, favoriteHandler *handler.FavoriteHandler) {
	favorites := rg.Group("/favorites")
	{
		favorites.POST("/movies", favoriteHandler.AddMovieToFavorites)
		favorites.DELETE("/movies/:movie_id", favoriteHandler.RemoveMovieFromFavorites)
		favorites.POST("/cinemas", favoriteHandler.AddCinemaToFavorites)
		favorites.DELETE("/cinemas/:cinema_id", favoriteHandler.RemoveCinemaFromFavorites)
		favorites.GET("", favoriteHandler.GetUserFavorites)
	}
}

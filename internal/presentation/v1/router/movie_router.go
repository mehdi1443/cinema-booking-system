package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpMovieRouter(r *gin.RouterGroup, movieHandler *handler.MovieHandler) {
	movies := r.Group("/movies")
	{
		movies.POST("", movieHandler.CreateMovie)
		movies.GET("", movieHandler.GetMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.PUT("/:id", movieHandler.UpdateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)

		movies.POST("/:id/poster", movieHandler.UploadPoster)
		movies.DELETE("/:id/poster", movieHandler.DeletePoster)

		movies.POST("/:id/trailer", movieHandler.UploadTrailer)
		movies.DELETE("/:id/trailer", movieHandler.DeleteTrailer)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpCinemaRouter(r *gin.RouterGroup, cinemaHandler *handler.CinemaHandler) {
	cinemas := r.Group("/cinemas")
	{
		cinemas.POST("", cinemaHandler.CreateCinema)
		cinemas.GET("", cinemaHandler.GetCinemas)
		cinemas.GET("/:id", cinemaHandler.GetCinemaByID)
		cinemas.PUT("/:id", cinemaHandler.UpdateCinema)
		cinemas.DELETE("/:id", cinemaHandler.DeleteCinema)

		cinemas.POST("/:id/image", cinemaHandler.UploadImage)
		cinemas.DELETE("/:id/image", cinemaHandler.DeleteImage)
	}
}

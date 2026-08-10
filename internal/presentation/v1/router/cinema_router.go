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

		cinemas.POST("/:id/image", cinemaHandler.UploadImage)
		cinemas.DELETE("/:id/image", cinemaHandler.DeleteImage)

		cinemas.POST("/:id/facilities/:facility_id", cinemaHandler.AddFacility)
		cinemas.DELETE("/:id/facilities/:facility_id", cinemaHandler.RemoveFacility)
		cinemas.PUT("/:id/facilities", cinemaHandler.UpdateFacilities)
	}
}

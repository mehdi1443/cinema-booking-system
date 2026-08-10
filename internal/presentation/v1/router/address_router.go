package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpAddressRouter(r *gin.RouterGroup, addressHandler *handler.AddressHandler) {
	// Province Endpoints
	provinces := r.Group("/provinces")
	{
		provinces.POST("", addressHandler.CreateProvince)
		provinces.GET("", addressHandler.GetAllProvinces)
		provinces.GET("/:id", addressHandler.GetProvinceByID)
		provinces.PUT("/:id", addressHandler.UpdateProvince)
		provinces.DELETE("/:id", addressHandler.DeleteProvince)
	}

	// City Endpoints
	cities := r.Group("/cities")
	{
		cities.POST("", addressHandler.CreateCity)
		cities.GET("", addressHandler.GetCities)
		cities.GET("/:id", addressHandler.GetCityByID)
		cities.PUT("/:id", addressHandler.UpdateCity)
		cities.DELETE("/:id", addressHandler.DeleteCity)
	}
}

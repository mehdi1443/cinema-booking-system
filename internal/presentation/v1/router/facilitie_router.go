package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpFacilityRouter(r *gin.RouterGroup, facilityHandler *handler.FacilityHandler) {
	facilities := r.Group("/facilities")
	{
		facilities.POST("", facilityHandler.CreateFacility)
		facilities.GET("", facilityHandler.GetAllFacilities)
		facilities.GET("/:id", facilityHandler.GetFacilityByID)
		facilities.PUT("/:id", facilityHandler.UpdateFacility)
		facilities.DELETE("/:id", facilityHandler.DeleteFacility)
	}
}

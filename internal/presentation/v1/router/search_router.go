package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpSearchRouter(
	r *gin.RouterGroup,
	searchHandler *handler.SearchHandler,
) {
	searchRoutes := r.Group("/search")
	{
		searchRoutes.GET("/", searchHandler.Search)
		searchRoutes.GET("/suggest", searchHandler.Suggest)
		searchRoutes.POST("/sync", searchHandler.SyncAll)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpGenreRouter(
	rg *gin.RouterGroup,
	h *handler.GenreHandler,
) {
	genres := rg.Group("/genres")
	{
		genres.POST("", h.CreateGenre)
		genres.GET("", h.GetGenres)
		genres.GET("/:id", h.GetGenreByID)
		genres.PUT("/:id", h.UpdateGenre)
		genres.DELETE("/:id", h.DeleteGenre)
	}
}

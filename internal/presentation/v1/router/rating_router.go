package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpRatingRouter(rg *gin.RouterGroup, ratingHandler *handler.RatingHandler) {
	rg.GET("/movies/:id/ratings", ratingHandler.GetMovieRatings)
	rg.GET("/cinemas/:id/ratings", ratingHandler.GetCinemaRatings)

	rg.POST("/ratings", ratingHandler.SubmitRating)
	rg.GET("/movies/:id/user-rating", ratingHandler.GetUserMovieRating)
	rg.GET("/cinemas/:id/user-rating", ratingHandler.GetUserCinemaRating)
	rg.DELETE("/ratings/:id", ratingHandler.DeleteRating)
}

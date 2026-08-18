package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type RatingHandler struct {
	ratingService service_contract.RatingServiceInterface
}

func NewRatingHandler(ratingService service_contract.RatingServiceInterface) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

// POST /ratings
func (h *RatingHandler) SubmitRating(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	var req dto.SubmitRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	rating, err := h.ratingService.SubmitRating(c.Request.Context(), userID, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] SubmitRating failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Exactly one target (movie_id or cinema_id) must be specified", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Rating submitted successfully", rating)
}

// GET /movies/:id/user-rating
func (h *RatingHandler) GetUserMovieRating(c *gin.Context) {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	rating, err := h.ratingService.GetUserRatingForMovie(c.Request.Context(), userID, movieID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetUserMovieRating failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User movie rating retrieved successfully", rating)
}

// GET /cinemas/:id/user-rating
func (h *RatingHandler) GetUserCinemaRating(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	rating, err := h.ratingService.GetUserRatingForCinema(c.Request.Context(), userID, cinemaID)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetUserCinemaRating failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User cinema rating retrieved successfully", rating)
}

// GET /movies/:id/ratings
func (h *RatingHandler) GetMovieRatings(c *gin.Context) {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	ratings, err := h.ratingService.GetByMovieID(c.Request.Context(), movieID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetMovieRatings failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie ratings retrieved successfully", ratings)
}

// GET /cinemas/:id/ratings
func (h *RatingHandler) GetCinemaRatings(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	ratings, err := h.ratingService.GetByCinemaID(c.Request.Context(), cinemaID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemaRatings failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema ratings retrieved successfully", ratings)
}

// DELETE /ratings/:id
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	ratingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid rating ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	if err := h.ratingService.Delete(c.Request.Context(), userID, ratingID); err != nil {
		log.Printf("[SERVICE ERROR] DeleteRating failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrRatingNotFound):
			response.Error(c, http.StatusNotFound, "Rating not found", err)
		case errors.Is(err, exception.ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "You are not allowed to delete this rating", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Rating deleted successfully", nil)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type RatingHandler struct {
	ratingService service_contract.RatingServiceInterface
	log           domainlogger.Logger
}

func NewRatingHandler(ratingService service_contract.RatingServiceInterface, log domainlogger.Logger) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
		log:           log,
	}
}

// -----------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------

// POST /ratings
func (h *RatingHandler) SubmitRating(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	var req dto.SubmitRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	rating, err := h.ratingService.SubmitRating(c.Request.Context(), userID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Rating submitted successfully", rating)
}

// GET /movies/:id/user-rating
func (h *RatingHandler) GetUserMovieRating(c *gin.Context) {
	movieID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	rating, err := h.ratingService.GetUserRatingForMovie(c.Request.Context(), userID, movieID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User movie rating retrieved successfully", rating)
}

// GET /cinemas/:id/user-rating
func (h *RatingHandler) GetUserCinemaRating(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	rating, err := h.ratingService.GetUserRatingForCinema(c.Request.Context(), userID, cinemaID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "User cinema rating retrieved successfully", rating)
}

// GET /movies/:id/ratings
func (h *RatingHandler) GetMovieRatings(c *gin.Context) {
	movieID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	ratings, err := h.ratingService.GetByMovieID(c.Request.Context(), movieID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Movie ratings retrieved successfully", ratings)
}

// GET /cinemas/:id/ratings
func (h *RatingHandler) GetCinemaRatings(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	ratings, err := h.ratingService.GetByCinemaID(c.Request.Context(), cinemaID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema ratings retrieved successfully", ratings)
}

// DELETE /ratings/:id
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	ratingID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	if err := h.ratingService.Delete(c.Request.Context(), userID, ratingID); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Rating deleted successfully", nil)
}

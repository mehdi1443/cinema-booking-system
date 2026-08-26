package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type FavoriteHandler struct {
	favoriteService     service_contract.FavoriteServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 domainlogger.Logger
}

func NewFavoriteHandler(
	favoriteService service_contract.FavoriteServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log domainlogger.Logger,
) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService:     favoriteService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// -----------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------

// POST /favorites/movies
func (h *FavoriteHandler) AddMovieToFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	var req dto.AddFavoriteMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	favorite, err := h.favoriteService.AddMovieToFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Movie added to favorites successfully", favorite)
}

// DELETE /favorites/movies/:movie_id
func (h *FavoriteHandler) RemoveMovieFromFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	movieID, ok := parseParamID(c, "movie_id")
	if !ok {
		return
	}

	if err := h.favoriteService.RemoveMovieFromFavorites(c.Request.Context(), tokenData.UserID, movieID); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Movie removed from favorites successfully", nil)
}

// POST /favorites/cinemas
func (h *FavoriteHandler) AddCinemaToFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	var req dto.AddFavoriteCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	favorite, err := h.favoriteService.AddCinemaToFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Cinema added to favorites successfully", favorite)
}

// DELETE /favorites/cinemas/:cinema_id
func (h *FavoriteHandler) RemoveCinemaFromFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	cinemaID, ok := parseParamID(c, "cinema_id")
	if !ok {
		return
	}

	if err := h.favoriteService.RemoveCinemaFromFavorites(c.Request.Context(), tokenData.UserID, cinemaID); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema removed from favorites successfully", nil)
}

// GET /favorites
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	var req dto.FavoriteFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.favoriteService.GetUserFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Favorites retrieved successfully", result)
}

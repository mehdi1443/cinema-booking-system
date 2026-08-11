package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type FavoriteHandler struct {
	favoriteService     service_contract.FavoriteServiceInterface
	contexttokenService service_contract.ContextTokenServiceInterface
	log                 logger.Logger
}

func NewFavoriteHandler(
	favoriteService service_contract.FavoriteServiceInterface,
	contexttokenService service_contract.ContextTokenServiceInterface,
	log logger.Logger,
) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService:     favoriteService,
		contexttokenService: contexttokenService,
		log:                 log,
	}
}

// POST /favorites/movies
func (h *FavoriteHandler) AddMovieToFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "AddMovieToFavorites missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error(), err)
		return
	}

	var req dto.AddFavoriteMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for AddMovieToFavorites", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	favorite, err := h.favoriteService.AddMovieToFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrDuplicateKey):
			h.log.Warn(logger.Validation, logger.Api, "AddMovieToFavorites failed: duplicate record", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "Movie already in favorites", err)

		case errors.Is(err, exception.ErrMovieNotFound):
			h.log.Warn(logger.Validation, logger.Api, "AddMovieToFavorites failed: movie not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Movie not found", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Insert, "AddMovieToFavorites database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "AddMovieToFavorites unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Movie added to favorites successfully", favorite)
}

// DELETE /favorites/movies/:movie_id
func (h *FavoriteHandler) RemoveMovieFromFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "RemoveMovieFromFavorites missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error(), err)
		return
	}

	movieID, err := strconv.ParseInt(c.Param("movie_id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid movie ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	if err := h.favoriteService.RemoveMovieFromFavorites(c.Request.Context(), tokenData.UserID, movieID); err != nil {
		switch {
		case errors.Is(err, exception.ErrRecordNotFound):
			h.log.Warn(logger.Validation, logger.Api, "RemoveMovieFromFavorites failed: record not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Movie not found in favorites", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Delete, "RemoveMovieFromFavorites database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "RemoveMovieFromFavorites unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie removed from favorites successfully", nil)
}

// POST /favorites/cinemas
func (h *FavoriteHandler) AddCinemaToFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "AddCinemaToFavorites missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error(), err)
		return
	}

	var req dto.AddFavoriteCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid payload for AddCinemaToFavorites", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	favorite, err := h.favoriteService.AddCinemaToFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrDuplicateKey):
			h.log.Warn(logger.Validation, logger.Api, "AddCinemaToFavorites failed: duplicate record", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusConflict, "Cinema already in favorites", err)

		case errors.Is(err, exception.ErrCinemaNotFound):
			h.log.Warn(logger.Validation, logger.Api, "AddCinemaToFavorites failed: cinema not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusBadRequest, "Cinema not found", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Insert, "AddCinemaToFavorites database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "AddCinemaToFavorites unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Cinema added to favorites successfully", favorite)
}

// DELETE /favorites/cinemas/:cinema_id
func (h *FavoriteHandler) RemoveCinemaFromFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "RemoveCinemaFromFavorites missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error(), err)
		return
	}

	cinemaID, err := strconv.ParseInt(c.Param("cinema_id"), 10, 64)
	if err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid cinema ID parameter", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	if err := h.favoriteService.RemoveCinemaFromFavorites(c.Request.Context(), tokenData.UserID, cinemaID); err != nil {
		switch {
		case errors.Is(err, exception.ErrRecordNotFound):
			h.log.Warn(logger.Validation, logger.Api, "RemoveCinemaFromFavorites failed: record not found", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusNotFound, "Cinema not found in favorites", err)

		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Delete, "RemoveCinemaFromFavorites database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "RemoveCinemaFromFavorites unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema removed from favorites successfully", nil)
}

// GET /favorites
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	tokenData, err := h.contexttokenService.GetTokenFromGinContext(c)
	if err != nil {
		h.log.Warn(logger.Internal, logger.Api, "GetUserFavorites missing token context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusUnauthorized, exception.ErrGetTokenFromContext.Error(), err)
		return
	}

	var req dto.FavoriteFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.Warn(logger.Validation, logger.Api, "Invalid query parameters for GetUserFavorites", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	result, err := h.favoriteService.GetUserFavorites(c.Request.Context(), tokenData.UserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrDatabase):
			h.log.Error(logger.Postgres, logger.Select, "GetUserFavorites database error", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			h.log.Error(logger.Internal, logger.Api, "GetUserFavorites unexpected failure", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Favorites retrieved successfully", result)
}

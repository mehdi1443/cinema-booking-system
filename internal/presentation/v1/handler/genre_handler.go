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

type GenreHandler struct {
	genreService service_contract.GenreServiceInterface
}

func NewGenreHandler(genreService service_contract.GenreServiceInterface) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

func (h *GenreHandler) CreateGenre(c *gin.Context) {
	var req dto.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	genre, err := h.genreService.CreateGenre(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateGenre failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDuplicateKey), errors.Is(err, exception.ErrGenreAlreadyExists):
			response.Error(c, http.StatusConflict, "Genre already exists", err)

		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid genre name", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Genre created successfully", genre)
}

func (h *GenreHandler) GetGenreByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid genre ID", err)
		return
	}

	genre, err := h.genreService.GetGenreByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetGenreByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrGenreNotFound):
			response.Error(c, http.StatusNotFound, "Genre not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Genre retrieved successfully", genre)
}

func (h *GenreHandler) GetGenres(c *gin.Context) {
	var req dto.GenreFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filter parameters", err)
		return
	}

	genres, total, err := h.genreService.GetGenres(c.Request.Context(), req)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetGenres failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	} else if page > 1000 {
		page = 1000
	}

	limit := req.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	meta := response.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Genres retrieved successfully", genres, meta)
}

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid genre ID", err)
		return
	}

	var req dto.UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedGenre, err := h.genreService.UpdateGenre(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateGenre failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrGenreNotFound):
			response.Error(c, http.StatusNotFound, "Genre not found", err)

		case errors.Is(err, exception.ErrDuplicateKey), errors.Is(err, exception.ErrGenreAlreadyExists):
			response.Error(c, http.StatusConflict, "Genre name already exists", err)

		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid genre name", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Genre updated successfully", updatedGenre)
}

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid genre ID", err)
		return
	}

	if err := h.genreService.DeleteGenre(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteGenre failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrGenreNotFound):
			response.Error(c, http.StatusNotFound, "Genre not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Genre deleted successfully", nil)
}

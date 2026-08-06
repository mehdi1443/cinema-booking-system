package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type GenreHandler struct {
	genreService service_contract.GenreUseCaseInterface
}

func NewGenreHandler(genreService service_contract.GenreUseCaseInterface) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

// POST /genres
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
		case errors.Is(err, exception.ErrGenreAlreadyExists):
			response.Error(c, http.StatusConflict, "Genre already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Genre created successfully", genre)
}

// GET /genres/:id
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

// GET /genres
func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.genreService.GetAllGenres(c.Request.Context())
	if err != nil {
		log.Printf("[SERVICE ERROR] GetAllGenres failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Genres retrieved successfully", genres)
}

// PUT /genres/:id
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

		case errors.Is(err, exception.ErrGenreAlreadyExists):
			response.Error(c, http.StatusConflict, "Genre name already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Genre updated successfully", updatedGenre)
}

// DELETE /genres/:id
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

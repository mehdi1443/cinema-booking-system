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

const (
	maxPosterSize  = 5 * 1024 * 1024   // 5 MB
	maxTrailerSize = 100 * 1024 * 1024 // 100 MB
)

type MovieHandler struct {
	movieService service_contract.MovieServiceInterface
}

func NewMovieHandler(movieService service_contract.MovieServiceInterface) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

// POST /movies
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req dto.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	movie, err := h.movieService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateMovie failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid input data or enum value", err)

		case errors.Is(err, exception.ErrMovieAlreadyExists):
			response.Error(c, http.StatusConflict, "Movie already exists", err)

		case errors.Is(err, exception.ErrGenreNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified genres were not found", err)

		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified persons were not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Movie created successfully", movie)
}

// GET /movies/:id
func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	movie, err := h.movieService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetMovieByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie retrieved successfully", movie)
}

// GET /movies
func (h *MovieHandler) GetMovies(c *gin.Context) {
	var req dto.MovieFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filter parameters", err)
		return
	}

	req.Normalize()

	movies, total, err := h.movieService.GetMovies(c.Request.Context(), req)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetMovies failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid filter parameter or enum value", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	meta := response.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Movies retrieved successfully", movies, meta)
}

// PUT /movies/:id
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	var req dto.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedMovie, err := h.movieService.Update(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateMovie failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid input data or enum value", err)

		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrGenreNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified genres were not found", err)

		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified persons were not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie updated successfully", updatedMovie)
}

// DELETE /movies/:id
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	if err := h.movieService.Delete(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteMovie failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie deleted successfully", nil)
}

// POST /movies/:id/poster
func (h *MovieHandler) UploadPoster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	fileHeader, err := c.FormFile("poster")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Poster image file is required", err)
		return
	}

	if fileHeader.Size > maxPosterSize {
		response.Error(c, http.StatusRequestEntityTooLarge, "Poster size exceeds allowed limit (5MB)", exception.ErrFileSizeTooLarge)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read poster file", err)
		return
	}
	defer file.Close()

	movie, err := h.movieService.UploadPoster(c.Request.Context(), id, file)
	if err != nil {
		log.Printf("[SERVICE ERROR] UploadPoster failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrInvalidFileType):
			response.Error(c, http.StatusBadRequest, "Invalid poster image format", err)

		case errors.Is(err, exception.ErrFileSizeTooLarge):
			response.Error(c, http.StatusRequestEntityTooLarge, "Poster size exceeds allowed limit (5MB)", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to upload poster", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie poster uploaded successfully", movie)
}

// DELETE /movies/:id/poster
func (h *MovieHandler) DeletePoster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	if err := h.movieService.DeletePoster(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeletePoster failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to delete poster", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie poster deleted successfully", nil)
}

// POST /movies/:id/trailer
func (h *MovieHandler) UploadTrailer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	fileHeader, err := c.FormFile("trailer")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Trailer video file is required", err)
		return
	}

	if fileHeader.Size > maxTrailerSize {
		response.Error(c, http.StatusRequestEntityTooLarge, "Trailer size exceeds allowed limit (100MB)", exception.ErrFileSizeTooLarge)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read trailer file", err)
		return
	}
	defer file.Close()

	movie, err := h.movieService.UploadTrailer(c.Request.Context(), id, file)
	if err != nil {
		log.Printf("[SERVICE ERROR] UploadTrailer failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrInvalidFileType):
			response.Error(c, http.StatusBadRequest, "Invalid trailer video format", err)

		case errors.Is(err, exception.ErrFileSizeTooLarge):
			response.Error(c, http.StatusRequestEntityTooLarge, "Trailer size exceeds allowed limit (100MB)", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to upload trailer", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie trailer uploaded successfully", movie)
}

// DELETE /movies/:id/trailer
func (h *MovieHandler) DeleteTrailer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	if err := h.movieService.DeleteTrailer(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteTrailer failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrMovieNotFound):
			response.Error(c, http.StatusNotFound, "Movie not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to delete trailer", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie trailer deleted successfully", nil)
}

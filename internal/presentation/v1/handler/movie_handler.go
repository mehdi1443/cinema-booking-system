package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
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
// Support optional query params: /movies?status=... or /movies?genre_id=...
func (h *MovieHandler) GetMovies(c *gin.Context) {
	ctx := c.Request.Context()

	if statusParam := c.Query("status"); statusParam != "" {
		movies, err := h.movieService.GetByStatus(ctx, enum.MovieStatus(statusParam))
		if err != nil {
			log.Printf("[SERVICE ERROR] GetByStatus failed: %v", err)
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve movies by status", err)
			return
		}
		response.Success(c, http.StatusOK, "Movies retrieved successfully", movies)
		return
	}

	if genreIDParam := c.Query("genre_id"); genreIDParam != "" {
		genreID, err := strconv.ParseInt(genreIDParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid genre ID format", err)
			return
		}

		movies, err := h.movieService.GetByGenreID(ctx, genreID)
		if err != nil {
			log.Printf("[SERVICE ERROR] GetByGenreID failed: %v", err)
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve movies by genre", err)
			return
		}
		response.Success(c, http.StatusOK, "Movies retrieved successfully", movies)
		return
	}

	movies, err := h.movieService.GetAll(ctx)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetAllMovies failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movies retrieved successfully", movies)
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

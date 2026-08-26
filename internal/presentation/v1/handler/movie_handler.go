package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

const (
	maxPosterSize  = 5 * 1024 * 1024
	maxTrailerSize = 100 * 1024 * 1024
)

type MovieHandler struct {
	movieService service_contract.MovieServiceInterface
	log          domainlogger.Logger
}

func NewMovieHandler(
	movieService service_contract.MovieServiceInterface,
	log domainlogger.Logger,
) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
		log:          log,
	}
}

// -----------------------------------------------------------------------------
// POST /movies
// -----------------------------------------------------------------------------

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req dto.CreateMovieRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	movie, err := h.movieService.Create(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log) // <--- پاس دادن لاگر به ارور هندلر
		return
	}

	response.Success(c, http.StatusCreated, "Movie created successfully", movie)
}

// -----------------------------------------------------------------------------
// GET /movies/:id
// -----------------------------------------------------------------------------

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	movie, err := h.movieService.GetByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie retrieved successfully", movie)
}

// -----------------------------------------------------------------------------
// GET /movies
// -----------------------------------------------------------------------------

func (h *MovieHandler) GetMovies(c *gin.Context) {
	var req dto.MovieFilterRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid filter parameters")
		return
	}

	req.Normalize()

	movies, total, err := h.movieService.GetMovies(c.Request.Context(), req)
	if err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	meta := response.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Movies retrieved successfully", movies, meta)
}

// -----------------------------------------------------------------------------
// PUT /movies/:id
// -----------------------------------------------------------------------------

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateMovieRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedMovie, err := h.movieService.Update(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie updated successfully", updatedMovie)
}

// -----------------------------------------------------------------------------
// DELETE /movies/:id
// -----------------------------------------------------------------------------

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.movieService.Delete(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie deleted successfully", nil)
}

// -----------------------------------------------------------------------------
// POST /movies/:id/poster
// -----------------------------------------------------------------------------

func (h *MovieHandler) UploadPoster(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("poster")
	if err != nil {
		HandleBadRequest(c, "Poster image file is required")
		return
	}

	if fileHeader.Size > maxPosterSize {
		Handle(c, exception.ErrFileSizeTooLarge, h.log) // <---
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		Handle(c, exception.ErrInternal, h.log) // <---
		return
	}
	defer file.Close()

	movie, err := h.movieService.UploadPoster(c.Request.Context(), id, file)
	if err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie poster uploaded successfully", movie)
}

// -----------------------------------------------------------------------------
// DELETE /movies/:id/poster
// -----------------------------------------------------------------------------

func (h *MovieHandler) DeletePoster(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.movieService.DeletePoster(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie poster deleted successfully", nil)
}

// -----------------------------------------------------------------------------
// POST /movies/:id/trailer
// -----------------------------------------------------------------------------

func (h *MovieHandler) UploadTrailer(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("trailer")
	if err != nil {
		HandleBadRequest(c, "Trailer video file is required")
		return
	}

	if fileHeader.Size > maxTrailerSize {
		Handle(c, exception.ErrFileSizeTooLarge, h.log) // <---
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		Handle(c, exception.ErrInternal, h.log) // <---
		return
	}
	defer file.Close()

	movie, err := h.movieService.UploadTrailer(c.Request.Context(), id, file)
	if err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie trailer uploaded successfully", movie)
}

// -----------------------------------------------------------------------------
// DELETE /movies/:id/trailer
// -----------------------------------------------------------------------------

func (h *MovieHandler) DeleteTrailer(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.movieService.DeleteTrailer(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log) // <---
		return
	}

	response.Success(c, http.StatusOK, "Movie trailer deleted successfully", nil)
}

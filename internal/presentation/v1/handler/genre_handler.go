package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
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

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func parseGenreID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		HandleBadRequest(c, "Invalid genre ID")
		return 0, false
	}
	return id, true
}

// -----------------------------------------------------------------------------
// POST /genres
// -----------------------------------------------------------------------------

func (h *GenreHandler) CreateGenre(c *gin.Context) {
	var req dto.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	genre, err := h.genreService.CreateGenre(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Genre created successfully", genre)
}

// -----------------------------------------------------------------------------
// GET /genres/:id
// -----------------------------------------------------------------------------

func (h *GenreHandler) GetGenreByID(c *gin.Context) {
	id, ok := parseGenreID(c)
	if !ok {
		return
	}

	genre, err := h.genreService.GetGenreByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Genre retrieved successfully", genre)
}

// -----------------------------------------------------------------------------
// GET /genres
// -----------------------------------------------------------------------------

func (h *GenreHandler) GetGenres(c *gin.Context) {
	var req dto.GenreFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid filter parameters")
		return
	}

	req.Normalize()

	genres, total, err := h.genreService.GetGenres(c.Request.Context(), req)
	if err != nil {
		Handle(c, err)
		return
	}

	meta := response.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Genres retrieved successfully", genres, meta)
}

// -----------------------------------------------------------------------------
// PUT /genres/:id
// -----------------------------------------------------------------------------

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	id, ok := parseGenreID(c)
	if !ok {
		return
	}

	var req dto.UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedGenre, err := h.genreService.UpdateGenre(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Genre updated successfully", updatedGenre)
}

// -----------------------------------------------------------------------------
// DELETE /genres/:id
// -----------------------------------------------------------------------------

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	id, ok := parseGenreID(c)
	if !ok {
		return
	}

	if err := h.genreService.DeleteGenre(c.Request.Context(), id); err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Genre deleted successfully", nil)
}

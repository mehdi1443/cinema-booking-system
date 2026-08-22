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

const maxCinemaImageSize = 5 * 1024 * 1024 // 5 MB

type CinemaHandler struct {
	cinemaService service_contract.CinemaServiceInterface
}

func NewCinemaHandler(cinemaService service_contract.CinemaServiceInterface) *CinemaHandler {
	return &CinemaHandler{
		cinemaService: cinemaService,
	}
}

// POST /cinemas
func (h *CinemaHandler) CreateCinema(c *gin.Context) {
	var req dto.CreateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	cinema, err := h.cinemaService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateCinema failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid input data or enum value", err)

		case errors.Is(err, exception.ErrDuplicateKey) || errors.Is(err, exception.ErrCinemaAlreadyExists):
			response.Error(c, http.StatusConflict, "Cinema already exists", err)

		case errors.Is(err, exception.ErrCityNotFound):
			response.Error(c, http.StatusBadRequest, "Specified city was not found", err)

		case errors.Is(err, exception.ErrFacilityNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified facilities were not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Cinema created successfully", cinema)
}

// GET /cinemas/:id
func (h *CinemaHandler) GetCinemaByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	cinema, err := h.cinemaService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemaByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCinemaNotFound):
			response.Error(c, http.StatusNotFound, "Cinema not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema retrieved successfully", cinema)
}

// GET /cinemas
func (h *CinemaHandler) GetCinemas(c *gin.Context) {
	var req dto.CinemaFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filter parameters", err)
		return
	}

	req.Normalize() // از PaginationRequest به ارث رسیده

	cinemas, total, err := h.cinemaService.GetCinemas(c.Request.Context(), req)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemas failed: %v", err)
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

	response.SuccessWithMeta(c, http.StatusOK, "Cinemas retrieved successfully", cinemas, meta)
}

// PUT /cinemas/:id
func (h *CinemaHandler) UpdateCinema(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	var req dto.UpdateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedCinema, err := h.cinemaService.Update(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateCinema failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "Invalid input data or enum value", err)

		case errors.Is(err, exception.ErrCinemaNotFound):
			response.Error(c, http.StatusNotFound, "Cinema not found", err)

		case errors.Is(err, exception.ErrCityNotFound):
			response.Error(c, http.StatusBadRequest, "Specified city was not found", err)

		case errors.Is(err, exception.ErrFacilityNotFound):
			response.Error(c, http.StatusBadRequest, "One or more specified facilities were not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema updated successfully", updatedCinema)
}

// DELETE /cinemas/:id
func (h *CinemaHandler) DeleteCinema(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	if err := h.cinemaService.Delete(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteCinema failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCinemaNotFound):
			response.Error(c, http.StatusNotFound, "Cinema not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema deleted successfully", nil)
}

// POST /cinemas/:id/image
func (h *CinemaHandler) UploadImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Image file is required", err)
		return
	}

	if fileHeader.Size > maxCinemaImageSize {
		response.Error(c, http.StatusRequestEntityTooLarge, "Image size exceeds allowed limit (5MB)", exception.ErrFileSizeTooLarge)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read image file", err)
		return
	}
	defer file.Close()

	cinema, err := h.cinemaService.UploadImage(c.Request.Context(), id, file, fileHeader.Filename)
	if err != nil {
		log.Printf("[SERVICE ERROR] UploadImage failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCinemaNotFound):
			response.Error(c, http.StatusNotFound, "Cinema not found", err)

		case errors.Is(err, exception.ErrInvalidFileType):
			response.Error(c, http.StatusBadRequest, "Invalid image format", err)

		case errors.Is(err, exception.ErrFileSizeTooLarge):
			response.Error(c, http.StatusRequestEntityTooLarge, "Image size exceeds allowed limit (5MB)", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to upload image", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema image uploaded successfully", cinema)
}

// DELETE /cinemas/:id/image
func (h *CinemaHandler) DeleteImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	if err := h.cinemaService.DeleteImage(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteImage failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCinemaNotFound):
			response.Error(c, http.StatusNotFound, "Cinema not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to delete image", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema image deleted successfully", nil)
}

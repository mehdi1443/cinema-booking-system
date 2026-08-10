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
		case errors.Is(err, exception.ErrCinemaAlreadyExists):
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
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	result, err := h.cinemaService.GetWithFilters(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemas failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinemas retrieved successfully", result)
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

// POST /cinemas/:id/facilities/:facility_id
func (h *CinemaHandler) AddFacility(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid facility ID", err)
		return
	}

	if err := h.cinemaService.AddFacility(c.Request.Context(), cinemaID, facilityID); err != nil {
		log.Printf("[SERVICE ERROR] AddFacility failed: %v", err)
		response.Error(c, http.StatusInternalServerError, "Failed to add facility to cinema", err)
		return
	}

	response.Success(c, http.StatusOK, "Facility added to cinema successfully", nil)
}

// DELETE /cinemas/:id/facilities/:facility_id
func (h *CinemaHandler) RemoveFacility(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid facility ID", err)
		return
	}

	if err := h.cinemaService.RemoveFacility(c.Request.Context(), cinemaID, facilityID); err != nil {
		log.Printf("[SERVICE ERROR] RemoveFacility failed: %v", err)
		response.Error(c, http.StatusInternalServerError, "Failed to remove facility from cinema", err)
		return
	}

	response.Success(c, http.StatusOK, "Facility removed from cinema successfully", nil)
}

// PUT /cinemas/:id/facilities
func (h *CinemaHandler) UpdateFacilities(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	var req struct {
		FacilityIDs []int64 `json:"facility_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.cinemaService.UpdateFacilities(c.Request.Context(), cinemaID, req.FacilityIDs); err != nil {
		log.Printf("[SERVICE ERROR] UpdateFacilities failed: %v", err)
		response.Error(c, http.StatusInternalServerError, "Failed to update cinema facilities", err)
		return
	}

	response.Success(c, http.StatusOK, "Cinema facilities updated successfully", nil)
}

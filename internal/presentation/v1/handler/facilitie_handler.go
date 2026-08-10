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

type FacilityHandler struct {
	facilityService service_contract.FacilityServiceInterface
}

func NewFacilityHandler(facilityService service_contract.FacilityServiceInterface) *FacilityHandler {
	return &FacilityHandler{
		facilityService: facilityService,
	}
}

// POST /facilities
func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req dto.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	facility, err := h.facilityService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateFacility failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrFacilityAlreadyExists):
			response.Error(c, http.StatusConflict, "Facility already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Facility created successfully", facility)
}

// GET /facilities/:id
func (h *FacilityHandler) GetFacilityByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid facility ID", err)
		return
	}

	facility, err := h.facilityService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetFacilityByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrFacilityNotFound):
			response.Error(c, http.StatusNotFound, "Facility not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Facility retrieved successfully", facility)
}

// GET /facilities
func (h *FacilityHandler) GetAllFacilities(c *gin.Context) {
	facilities, err := h.facilityService.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("[SERVICE ERROR] GetAllFacilities failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Facilities retrieved successfully", facilities)
}

// PUT /facilities/:id
func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid facility ID", err)
		return
	}

	var req dto.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedFacility, err := h.facilityService.Update(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateFacility failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrFacilityNotFound):
			response.Error(c, http.StatusNotFound, "Facility not found", err)

		case errors.Is(err, exception.ErrFacilityAlreadyExists):
			response.Error(c, http.StatusConflict, "Facility name already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Facility updated successfully", updatedFacility)
}

// DELETE /facilities/:id
func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid facility ID", err)
		return
	}

	if err := h.facilityService.Delete(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteFacility failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrFacilityNotFound):
			response.Error(c, http.StatusNotFound, "Facility not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Facility deleted successfully", nil)
}

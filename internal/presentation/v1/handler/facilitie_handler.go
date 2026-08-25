package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
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

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func parseFacilityID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		HandleBadRequest(c, "Invalid facility ID")
		return 0, false
	}
	return id, true
}

// -----------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------

func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req dto.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	facility, err := h.facilityService.Create(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Facility created successfully", facility)
}

func (h *FacilityHandler) GetFacilityByID(c *gin.Context) {
	id, ok := parseFacilityID(c)
	if !ok {
		return
	}

	facility, err := h.facilityService.GetByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facility retrieved successfully", facility)
}

func (h *FacilityHandler) GetAllFacilities(c *gin.Context) {
	facilities, err := h.facilityService.GetAll(c.Request.Context())
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facilities retrieved successfully", facilities)
}

func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, ok := parseFacilityID(c)
	if !ok {
		return
	}

	var req dto.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedFacility, err := h.facilityService.Update(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facility updated successfully", updatedFacility)
}

func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id, ok := parseFacilityID(c)
	if !ok {
		return
	}

	if err := h.facilityService.Delete(c.Request.Context(), id); err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facility deleted successfully", nil)
}

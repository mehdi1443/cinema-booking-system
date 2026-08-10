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

type AddressHandler struct {
	addressService service_contract.AddressServiceInterface
}

func NewAddressHandler(addressService service_contract.AddressServiceInterface) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// ==================== Province Endpoints ====================

// POST /provinces
func (h *AddressHandler) CreateProvince(c *gin.Context) {
	var req dto.CreateProvinceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	province, err := h.addressService.CreateProvince(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateProvince failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrProvinceAlreadyExists):
			response.Error(c, http.StatusConflict, "Province already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Province created successfully", province)
}

// GET /provinces/:id
func (h *AddressHandler) GetProvinceByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid province ID", err)
		return
	}

	province, err := h.addressService.GetProvinceByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetProvinceByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrProvinceNotFound):
			response.Error(c, http.StatusNotFound, "Province not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Province retrieved successfully", province)
}

// GET /provinces
func (h *AddressHandler) GetAllProvinces(c *gin.Context) {
	provinces, err := h.addressService.GetAllProvinces(c.Request.Context())
	if err != nil {
		log.Printf("[SERVICE ERROR] GetAllProvinces failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Provinces retrieved successfully", provinces)
}

// PUT /provinces/:id
func (h *AddressHandler) UpdateProvince(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid province ID", err)
		return
	}

	var req dto.UpdateProvinceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedProvince, err := h.addressService.UpdateProvince(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateProvince failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrProvinceNotFound):
			response.Error(c, http.StatusNotFound, "Province not found", err)

		case errors.Is(err, exception.ErrProvinceAlreadyExists):
			response.Error(c, http.StatusConflict, "Province name already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Province updated successfully", updatedProvince)
}

// DELETE /provinces/:id
func (h *AddressHandler) DeleteProvince(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid province ID", err)
		return
	}

	if err := h.addressService.DeleteProvince(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteProvince failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrProvinceNotFound):
			response.Error(c, http.StatusNotFound, "Province not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Province deleted successfully", nil)
}

// ==================== City Endpoints ====================

// POST /cities
func (h *AddressHandler) CreateCity(c *gin.Context) {
	var req dto.CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	city, err := h.addressService.CreateCity(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateCity failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCityAlreadyExists):
			response.Error(c, http.StatusConflict, "City already exists in this province", err)

		case errors.Is(err, exception.ErrProvinceNotFound):
			response.Error(c, http.StatusBadRequest, "Specified province was not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "City created successfully", city)
}

// GET /cities/:id
func (h *AddressHandler) GetCityByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid city ID", err)
		return
	}

	city, err := h.addressService.GetCityByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCityByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCityNotFound):
			response.Error(c, http.StatusNotFound, "City not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "City retrieved successfully", city)
}

// GET /cities
// Supports optional query param: /cities?province_id=...
func (h *AddressHandler) GetCities(c *gin.Context) {
	ctx := c.Request.Context()

	if provinceIDParam := c.Query("province_id"); provinceIDParam != "" {
		provinceID, err := strconv.ParseInt(provinceIDParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid province ID format", err)
			return
		}

		cities, err := h.addressService.GetCitiesByProvinceID(ctx, provinceID)
		if err != nil {
			log.Printf("[SERVICE ERROR] GetCitiesByProvinceID failed: %v", err)
			response.Error(c, http.StatusInternalServerError, "Failed to retrieve cities for specified province", err)
			return
		}

		response.Success(c, http.StatusOK, "Cities retrieved successfully", cities)
		return
	}

	cities, err := h.addressService.GetAllCities(ctx)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetAllCities failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cities retrieved successfully", cities)
}

// PUT /cities/:id
func (h *AddressHandler) UpdateCity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid city ID", err)
		return
	}

	var req dto.UpdateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedCity, err := h.addressService.UpdateCity(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateCity failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCityNotFound):
			response.Error(c, http.StatusNotFound, "City not found", err)

		case errors.Is(err, exception.ErrProvinceNotFound):
			response.Error(c, http.StatusBadRequest, "Specified province was not found", err)

		case errors.Is(err, exception.ErrCityAlreadyExists):
			response.Error(c, http.StatusConflict, "City name already exists in target province", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "City updated successfully", updatedCity)
}

// DELETE /cities/:id
func (h *AddressHandler) DeleteCity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid city ID", err)
		return
	}

	if err := h.addressService.DeleteCity(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteCity failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCityNotFound):
			response.Error(c, http.StatusNotFound, "City not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "City deleted successfully", nil)
}

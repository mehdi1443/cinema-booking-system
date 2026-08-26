package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type AddressHandler struct {
	addressService service_contract.AddressServiceInterface
	log            domainlogger.Logger
}

func NewAddressHandler(addressService service_contract.AddressServiceInterface, log domainlogger.Logger) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		log:            log,
	}
}

// ==================== Province Endpoints ====================

func (h *AddressHandler) CreateProvince(c *gin.Context) {
	var req dto.CreateProvinceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	province, err := h.addressService.CreateProvince(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Province created successfully", province)
}

func (h *AddressHandler) GetProvinceByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	province, err := h.addressService.GetProvinceByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Province retrieved successfully", province)
}

func (h *AddressHandler) GetAllProvinces(c *gin.Context) {
	provinces, err := h.addressService.GetAllProvinces(c.Request.Context())
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Provinces retrieved successfully", provinces)
}

func (h *AddressHandler) UpdateProvince(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateProvinceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedProvince, err := h.addressService.UpdateProvince(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Province updated successfully", updatedProvince)
}

func (h *AddressHandler) DeleteProvince(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.addressService.DeleteProvince(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Province deleted successfully", nil)
}

// ==================== City Endpoints ====================

func (h *AddressHandler) CreateCity(c *gin.Context) {
	var req dto.CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	city, err := h.addressService.CreateCity(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "City created successfully", city)
}

func (h *AddressHandler) GetCityByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	city, err := h.addressService.GetCityByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "City retrieved successfully", city)
}

func (h *AddressHandler) GetCities(c *gin.Context) {
	ctx := c.Request.Context()

	if provinceIDParam := c.Query("province_id"); provinceIDParam != "" {
		provinceID, err := strconv.ParseInt(provinceIDParam, 10, 64)
		if err != nil || provinceID <= 0 {
			HandleBadRequest(c, "Invalid province ID format")
			return
		}

		cities, err := h.addressService.GetCitiesByProvinceID(ctx, provinceID)
		if err != nil {
			Handle(c, err, h.log)
			return
		}

		response.Success(c, http.StatusOK, "Cities retrieved successfully", cities)
		return
	}

	cities, err := h.addressService.GetAllCities(ctx)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cities retrieved successfully", cities)
}

func (h *AddressHandler) UpdateCity(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedCity, err := h.addressService.UpdateCity(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "City updated successfully", updatedCity)
}

func (h *AddressHandler) DeleteCity(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.addressService.DeleteCity(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "City deleted successfully", nil)
}

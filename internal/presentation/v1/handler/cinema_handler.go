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

const maxCinemaImageSize = 5 * 1024 * 1024

type CinemaHandler struct {
	cinemaService service_contract.CinemaServiceInterface
	log           domainlogger.Logger
}

func NewCinemaHandler(cinemaService service_contract.CinemaServiceInterface, log domainlogger.Logger) *CinemaHandler {
	return &CinemaHandler{
		cinemaService: cinemaService,
		log:           log,
	}
}

func (h *CinemaHandler) SyncSearch(c *gin.Context) {
	if err := h.cinemaService.SyncAllToSearch(c.Request.Context()); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Background sync to Elasticsearch started successfully", nil)
}

func (h *CinemaHandler) CreateCinema(c *gin.Context) {
	var req dto.CreateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	cinema, err := h.cinemaService.Create(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Cinema created successfully", cinema)
}

func (h *CinemaHandler) GetCinemaByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	cinema, err := h.cinemaService.GetByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema retrieved successfully", cinema)
}

func (h *CinemaHandler) GetCinemas(c *gin.Context) {
	var req dto.CinemaFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid filter parameters")
		return
	}

	req.Normalize()

	cinemas, total, err := h.cinemaService.GetCinemas(c.Request.Context(), req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	meta := response.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Cinemas retrieved successfully", cinemas, meta)
}

func (h *CinemaHandler) UpdateCinema(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCinemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedCinema, err := h.cinemaService.Update(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema updated successfully", updatedCinema)
}

func (h *CinemaHandler) DeleteCinema(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.cinemaService.Delete(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema deleted successfully", nil)
}

func (h *CinemaHandler) UploadImage(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		HandleBadRequest(c, "Image file is required")
		return
	}

	if fileHeader.Size > maxCinemaImageSize {
		Handle(c, exception.ErrFileSizeTooLarge, h.log)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		Handle(c, exception.ErrInternal, h.log)
		return
	}
	defer file.Close()

	cinema, err := h.cinemaService.UploadImage(c.Request.Context(), id, file, fileHeader.Filename)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema image uploaded successfully", cinema)
}

func (h *CinemaHandler) DeleteImage(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.cinemaService.DeleteImage(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema image deleted successfully", nil)
}

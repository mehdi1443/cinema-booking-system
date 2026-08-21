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

type PersonHandler struct {
	personService service_contract.PersonServiceInterface
}

func NewPersonHandler(personService service_contract.PersonServiceInterface) *PersonHandler {
	return &PersonHandler{
		personService: personService,
	}
}

// POST /persons
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	person, err := h.personService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreatePerson failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonAlreadyExists):
			response.Error(c, http.StatusConflict, "Person already exists", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Person created successfully", person)
}

// GET /persons/:id
func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid person ID", err)
		return
	}

	person, err := h.personService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetPersonByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusNotFound, "Person not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Person retrieved successfully", person)
}

// GET /persons
func (h *PersonHandler) GetPersons(c *gin.Context) {
	var req dto.PersonFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid filter parameters", err)
		return
	}

	persons, total, err := h.personService.GetPersons(c.Request.Context(), req)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetPersons failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	meta := response.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Persons retrieved successfully", persons, meta)
}

// PUT /persons/:id
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid person ID", err)
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedPerson, err := h.personService.Update(c.Request.Context(), id, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdatePerson failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusNotFound, "Person not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Person updated successfully", updatedPerson)
}

// DELETE /persons/:id
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid person ID", err)
		return
	}

	if err := h.personService.Delete(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeletePerson failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusNotFound, "Person not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Person deleted successfully", nil)
}

// POST /persons/:id/image
func (h *PersonHandler) UploadImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid person ID", err)
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

	person, err := h.personService.UploadImage(c.Request.Context(), id, file)
	if err != nil {
		log.Printf("[SERVICE ERROR] UploadImage failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusNotFound, "Person not found", err)

		case errors.Is(err, exception.ErrInvalidFileType):
			response.Error(c, http.StatusBadRequest, "Invalid image format", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to upload image", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Person image uploaded successfully", person)
}

// DELETE /persons/:id/image
func (h *PersonHandler) DeleteImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid person ID", err)
		return
	}

	if err := h.personService.DeleteImage(c.Request.Context(), id); err != nil {
		log.Printf("[SERVICE ERROR] DeleteImage failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrPersonNotFound):
			response.Error(c, http.StatusNotFound, "Person not found", err)

		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)

		default:
			response.Error(c, http.StatusInternalServerError, "Failed to delete image", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Person image deleted successfully", nil)
}

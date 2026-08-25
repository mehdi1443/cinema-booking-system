package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

const maxPersonImageSize = 5 * 1024 * 1024 // 5 MB

type PersonHandler struct {
	personService service_contract.PersonServiceInterface
}

func NewPersonHandler(personService service_contract.PersonServiceInterface) *PersonHandler {
	return &PersonHandler{
		personService: personService,
	}
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func parsePersonID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		HandleBadRequest(c, "Invalid person ID")
		return 0, false
	}
	return id, true
}

// -----------------------------------------------------------------------------
// POST /persons
// -----------------------------------------------------------------------------

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	person, err := h.personService.Create(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Person created successfully", person)
}

// -----------------------------------------------------------------------------
// GET /persons/:id
// -----------------------------------------------------------------------------

func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	id, ok := parsePersonID(c)
	if !ok {
		return
	}

	person, err := h.personService.GetByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Person retrieved successfully", person)
}

// -----------------------------------------------------------------------------
// GET /persons
// -----------------------------------------------------------------------------

func (h *PersonHandler) GetPersons(c *gin.Context) {
	var req dto.PersonFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid filter parameters")
		return
	}

	req.Normalize()

	persons, total, err := h.personService.GetPersons(c.Request.Context(), req)
	if err != nil {
		Handle(c, err)
		return
	}

	meta := response.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	response.SuccessWithMeta(c, http.StatusOK, "Persons retrieved successfully", persons, meta)
}

// -----------------------------------------------------------------------------
// PUT /persons/:id
// -----------------------------------------------------------------------------

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id, ok := parsePersonID(c)
	if !ok {
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedPerson, err := h.personService.Update(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Person updated successfully", updatedPerson)
}

// -----------------------------------------------------------------------------
// DELETE /persons/:id
// -----------------------------------------------------------------------------

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id, ok := parsePersonID(c)
	if !ok {
		return
	}

	if err := h.personService.Delete(c.Request.Context(), id); err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Person deleted successfully", nil)
}

// -----------------------------------------------------------------------------
// POST /persons/:id/image
// -----------------------------------------------------------------------------

func (h *PersonHandler) UploadImage(c *gin.Context) {
	id, ok := parsePersonID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		HandleBadRequest(c, "Image file is required")
		return
	}

	if fileHeader.Size > maxPersonImageSize {
		Handle(c, exception.ErrFileSizeTooLarge)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		Handle(c, exception.ErrInternal)
		return
	}
	defer file.Close()

	person, err := h.personService.UploadImage(c.Request.Context(), id, file)
	if err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Person image uploaded successfully", person)
}

// -----------------------------------------------------------------------------
// DELETE /persons/:id/image
// -----------------------------------------------------------------------------

func (h *PersonHandler) DeleteImage(c *gin.Context) {
	id, ok := parsePersonID(c)
	if !ok {
		return
	}

	if err := h.personService.DeleteImage(c.Request.Context(), id); err != nil {
		Handle(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Person image deleted successfully", nil)
}

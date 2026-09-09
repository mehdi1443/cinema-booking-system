package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type TicketHandler struct {
	tickets service_contract.TicketServiceInterface
	log     domainlogger.Logger
}

func NewTicketHandler(tickets service_contract.TicketServiceInterface, log domainlogger.Logger) *TicketHandler {
	return &TicketHandler{tickets: tickets, log: log}
}

func (h *TicketHandler) ListMine(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	page, limit := parsePagination(c)
	result, total, err := h.tickets.ListMine(c.Request.Context(), userID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "Tickets retrieved successfully", result, response.Meta{Page: page, Limit: limit, Total: int(total)})
}

func (h *TicketHandler) Get(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.tickets.Get(c.Request.Context(), id, userID, hasPermission(c, enum.TicketsReadAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Ticket retrieved successfully", result)
}

func (h *TicketHandler) QRCode(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	png, err := h.tickets.QRCode(c.Request.Context(), id, userID, hasPermission(c, enum.TicketsReadAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	c.Header("Cache-Control", "private, no-store")
	c.Data(http.StatusOK, "image/png", png)
}

func (h *TicketHandler) PDF(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.tickets.PDF(c.Request.Context(), id, userID, hasPermission(c, enum.TicketsReadAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	c.Header("Cache-Control", "private, no-store")
	response.Success(c, http.StatusOK, "Ticket PDF URL generated successfully", result)
}

func (h *TicketHandler) Validate(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	var req dto.ValidateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.tickets.Validate(c.Request.Context(), req.Code, userID, hasPermission(c, enum.TicketsValidateAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Ticket validated successfully", result)
}

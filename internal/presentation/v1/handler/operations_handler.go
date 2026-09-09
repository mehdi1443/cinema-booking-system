package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type OperationsHandler struct {
	service service_contract.OperationsServiceInterface
	log     domainlogger.Logger
}

func NewOperationsHandler(service service_contract.OperationsServiceInterface, log domainlogger.Logger) *OperationsHandler {
	return &OperationsHandler{service: service, log: log}
}

func parseOptionalTime(c *gin.Context, key string, fallback time.Time) (time.Time, bool) {
	value := c.Query(key)
	if value == "" {
		return fallback, true
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		HandleBadRequest(c, key+" must be an RFC3339 timestamp")
		return time.Time{}, false
	}
	return parsed, true
}

func (h *OperationsHandler) SalesReport(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	now := time.Now()
	from, ok := parseOptionalTime(c, "from", now.Add(-30*24*time.Hour))
	if !ok {
		return
	}
	to, ok := parseOptionalTime(c, "to", now)
	if !ok {
		return
	}
	result, err := h.service.SalesReport(c.Request.Context(), cinemaID, from, to)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Cinema sales report retrieved successfully", result)
}

func (h *OperationsHandler) CreateSettlement(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	actorID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	var req dto.CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.service.CreateSettlement(c.Request.Context(), cinemaID, actorID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusCreated, "Settlement created successfully", result)
}

func (h *OperationsHandler) ListSettlements(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	page, limit := parsePagination(c)
	items, total, err := h.service.ListSettlements(c.Request.Context(), cinemaID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "Settlements retrieved successfully", items, response.Meta{Page: page, Limit: limit, Total: int(total)})
}

func (h *OperationsHandler) MarkSettlementPaid(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	actorID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	item, err := h.service.MarkSettlementPaid(c.Request.Context(), id, actorID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Settlement marked as paid", item)
}

func (h *OperationsHandler) ListAudit(c *gin.Context) {
	page, limit := parsePagination(c)
	filter := dto.OperationsFilterRequest{Page: page, Limit: limit}
	if value := c.Query("actor_user_id"); value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			HandleBadRequest(c, "Invalid actor_user_id")
			return
		}
		filter.ActorUserID = &id
	}
	if value := c.Query("from"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			HandleBadRequest(c, "from must be an RFC3339 timestamp")
			return
		}
		filter.From = &parsed
	}
	if value := c.Query("to"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			HandleBadRequest(c, "to must be an RFC3339 timestamp")
			return
		}
		filter.To = &parsed
	}
	items, total, err := h.service.ListAudit(c.Request.Context(), filter)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "Audit logs retrieved successfully", items, response.Meta{Page: page, Limit: limit, Total: int(total)})
}

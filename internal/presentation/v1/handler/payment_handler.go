package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type PaymentHandler struct {
	payments service_contract.PaymentServiceInterface
	log      domainlogger.Logger
}

func NewPaymentHandler(payments service_contract.PaymentServiceInterface, log domainlogger.Logger) *PaymentHandler {
	return &PaymentHandler{payments: payments, log: log}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	reservationID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	idempotencyKey := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if idempotencyKey == "" {
		HandleBadRequest(c, "Idempotency-Key header is required")
		return
	}
	result, created, err := h.payments.Create(c.Request.Context(), reservationID, userID, idempotencyKey)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	status := http.StatusOK
	message := "Payment retrieved successfully"
	if created {
		status = http.StatusCreated
		message = "Mock payment initialized successfully"
	}
	response.Success(c, status, message, result)
}

func (h *PaymentHandler) Get(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.payments.Get(c.Request.Context(), id, userID, hasPermission(c, enum.PaymentsReadAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Payment retrieved successfully", result)
}

func (h *PaymentHandler) ListForReservation(c *gin.Context) {
	reservationID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	page, limit := parsePagination(c)
	result, total, err := h.payments.ListForReservation(c.Request.Context(), reservationID, userID, hasPermission(c, enum.PaymentsReadAny), page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "Payments retrieved successfully", result, response.Meta{Page: page, Limit: limit, Total: int(total)})
}

func (h *PaymentHandler) Verify(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	var req dto.VerifyMockPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.payments.Verify(c.Request.Context(), id, userID, hasPermission(c, enum.PaymentsVerifyAny), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Mock payment processed successfully", result)
}

func (h *PaymentHandler) Refund(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.payments.Refund(c.Request.Context(), id, userID, hasPermission(c, enum.PaymentsRefundAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Mock payment refunded successfully", result)
}

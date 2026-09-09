package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpOperationsRouter(r *gin.RouterGroup, h *handler.OperationsHandler) {
	r.GET("/cinemas/:id/reports/sales", h.SalesReport)
	r.GET("/cinemas/:id/settlements", h.ListSettlements)
	r.POST("/cinemas/:id/settlements", h.CreateSettlement)
	r.PATCH("/settlements/:id/paid", h.MarkSettlementPaid)
	r.GET("/audit-logs", h.ListAudit)
}

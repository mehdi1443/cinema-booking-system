package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpPaymentRouter(r *gin.RouterGroup, h *handler.PaymentHandler) {
	r.POST("/reservations/:id/payments", h.Create)
	r.GET("/reservations/:id/payments", h.ListForReservation)
	r.GET("/payments/:id", h.Get)
	r.POST("/payments/:id/verify", h.Verify)
	r.POST("/payments/:id/refund", h.Refund)
}

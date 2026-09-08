package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpTicketRouter(r *gin.RouterGroup, h *handler.TicketHandler) {
	r.GET("/tickets", h.ListMine)
	r.GET("/tickets/:id", h.Get)
	r.GET("/tickets/:id/qr", h.QRCode)
	r.POST("/tickets/validate", h.Validate)
}

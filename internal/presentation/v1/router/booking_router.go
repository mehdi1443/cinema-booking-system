package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpBookingRouter(r *gin.RouterGroup, h *handler.BookingHandler) {
	r.POST("/cinemas/:id/halls", h.CreateHall)
	r.GET("/cinemas/:id/halls", h.ListHalls)

	r.GET("/halls/:id", h.GetHall)
	r.PUT("/halls/:id", h.UpdateHall)
	r.DELETE("/halls/:id", h.DeleteHall)
	r.POST("/halls/:id/seats", h.CreateSeats)
	r.GET("/halls/:id/seats", h.ListSeats)
	r.POST("/halls/:id/showtimes", h.CreateShowtime)

	r.PUT("/seats/:id", h.UpdateSeat)
	r.DELETE("/seats/:id", h.DeleteSeat)

	r.GET("/showtimes", h.ListShowtimes)
	r.GET("/showtimes/:id", h.GetShowtime)
	r.PUT("/showtimes/:id", h.UpdateShowtime)
	r.DELETE("/showtimes/:id", h.DeleteShowtime)
	r.PUT("/showtimes/:id/prices", h.UpdateShowtimePrices)
	r.GET("/showtimes/:id/seats", h.ListShowtimeSeats)

	r.POST("/reservations", h.CreateReservation)
	r.GET("/reservations", h.ListMyReservations)
	r.GET("/reservations/:id", h.GetReservation)
	r.DELETE("/reservations/:id", h.CancelReservation)
}

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

type BookingHandler struct {
	halls        service_contract.HallServiceInterface
	showtimes    service_contract.ShowtimeServiceInterface
	reservations service_contract.ReservationServiceInterface
	log          domainlogger.Logger
}

func NewBookingHandler(halls service_contract.HallServiceInterface, showtimes service_contract.ShowtimeServiceInterface, reservations service_contract.ReservationServiceInterface, log domainlogger.Logger) *BookingHandler {
	return &BookingHandler{halls: halls, showtimes: showtimes, reservations: reservations, log: log}
}

func (h *BookingHandler) CreateHall(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.CreateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.halls.CreateHall(c.Request.Context(), cinemaID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusCreated, "Hall created successfully", result)
}

func (h *BookingHandler) GetHall(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	result, err := h.halls.GetHall(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Hall retrieved successfully", result)
}

func (h *BookingHandler) ListHalls(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	result, err := h.halls.ListHalls(c.Request.Context(), cinemaID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Halls retrieved successfully", result)
}

func (h *BookingHandler) UpdateHall(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.halls.UpdateHall(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Hall updated successfully", result)
}

func (h *BookingHandler) DeleteHall(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	if err := h.halls.DeleteHall(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Hall deleted successfully", nil)
}

func (h *BookingHandler) CreateSeats(c *gin.Context) {
	hallID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.CreateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.halls.CreateSeats(c.Request.Context(), hallID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusCreated, "Seats created successfully", result)
}

func (h *BookingHandler) ListSeats(c *gin.Context) {
	hallID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	result, err := h.halls.ListSeats(c.Request.Context(), hallID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Seats retrieved successfully", result)
}

func (h *BookingHandler) UpdateSeat(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.halls.UpdateSeat(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Seat updated successfully", result)
}

func (h *BookingHandler) DeleteSeat(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	if err := h.halls.DeleteSeat(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Seat deleted successfully", nil)
}

func (h *BookingHandler) CreateShowtime(c *gin.Context) {
	hallID, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.CreateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.showtimes.Create(c.Request.Context(), hallID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusCreated, "Showtime created successfully", result)
}

func (h *BookingHandler) GetShowtime(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	result, err := h.showtimes.Get(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Showtime retrieved successfully", result)
}

func (h *BookingHandler) ListShowtimes(c *gin.Context) {
	var req dto.ShowtimeFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		HandleBadRequest(c, "Invalid filter parameters")
		return
	}
	result, total, err := h.showtimes.List(c.Request.Context(), &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.SuccessWithMeta(c, http.StatusOK, "Showtimes retrieved successfully", result, response.Meta{Page: req.Page, Limit: req.Limit, Total: int(total)})
}

func (h *BookingHandler) UpdateShowtime(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.showtimes.Update(c.Request.Context(), id, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Showtime updated successfully", result)
}

func (h *BookingHandler) DeleteShowtime(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	if err := h.showtimes.Delete(c.Request.Context(), id); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Showtime deleted successfully", nil)
}

func (h *BookingHandler) UpdateShowtimePrices(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateShowtimePricesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	if err := h.showtimes.UpdatePrices(c.Request.Context(), id, &req); err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Showtime prices updated successfully", nil)
}

func (h *BookingHandler) ListShowtimeSeats(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	result, err := h.showtimes.ListSeats(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Showtime seats retrieved successfully", result)
}

func (h *BookingHandler) CreateReservation(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}
	result, err := h.reservations.Create(c.Request.Context(), userID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusCreated, "Seats held successfully", result)
}

func (h *BookingHandler) GetReservation(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.reservations.Get(c.Request.Context(), id, userID, hasPermission(c, enum.ReservationsReadAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Reservation retrieved successfully", result)
}

func (h *BookingHandler) ListMyReservations(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.reservations.ListMine(c.Request.Context(), userID)
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Reservations retrieved successfully", result)
}

func (h *BookingHandler) CancelReservation(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}
	result, err := h.reservations.Cancel(c.Request.Context(), id, userID, hasPermission(c, enum.ReservationsCancelAny))
	if err != nil {
		Handle(c, err, h.log)
		return
	}
	response.Success(c, http.StatusOK, "Reservation cancelled successfully", result)
}

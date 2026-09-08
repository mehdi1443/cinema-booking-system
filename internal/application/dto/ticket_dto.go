package dto

import (
	"fmt"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type ValidateTicketRequest struct {
	Code string `json:"code" binding:"required,notblank,max=100"`
}

type TicketResponse struct {
	ID            int64                     `json:"id"`
	ReservationID int64                     `json:"reservation_id"`
	PaymentID     int64                     `json:"payment_id"`
	UserID        int64                     `json:"user_id"`
	ShowtimeID    int64                     `json:"showtime_id"`
	Code          string                    `json:"code"`
	Status        enum.TicketStatus         `json:"status"`
	QREndpoint    string                    `json:"qr_endpoint"`
	CinemaID      int64                     `json:"cinema_id,omitempty"`
	CinemaName    string                    `json:"cinema_name,omitempty"`
	HallID        int64                     `json:"hall_id,omitempty"`
	HallName      string                    `json:"hall_name,omitempty"`
	MovieID       int64                     `json:"movie_id,omitempty"`
	MovieTitle    string                    `json:"movie_title,omitempty"`
	StartsAt      *time.Time                `json:"starts_at,omitempty"`
	EndsAt        *time.Time                `json:"ends_at,omitempty"`
	Seats         []ReservationItemResponse `json:"seats"`
	UsedAt        *time.Time                `json:"used_at,omitempty"`
	UsedByUserID  *int64                    `json:"used_by_user_id,omitempty"`
	CreatedAt     *time.Time                `json:"created_at,omitempty"`
	UpdatedAt     *time.Time                `json:"updated_at,omitempty"`
}

func ToTicketResponse(ticket *entity.Ticket) *TicketResponse {
	if ticket == nil {
		return nil
	}
	result := &TicketResponse{
		ID: ticket.ID, ReservationID: ticket.ReservationID, PaymentID: ticket.PaymentID,
		UserID: ticket.UserID, Code: ticket.Code, Status: ticket.Status,
		QREndpoint: fmt.Sprintf("/v1/tickets/%d/qr", ticket.ID),
		Seats:      make([]ReservationItemResponse, 0), UsedAt: ticket.UsedAt,
		UsedByUserID: ticket.UsedByUserID, CreatedAt: helper.TimeToPtr(ticket.CreatedAt),
		UpdatedAt: helper.TimeToPtr(ticket.UpdatedAt),
	}
	reservation := ticket.Reservation
	if reservation == nil {
		return result
	}
	result.ShowtimeID = reservation.ShowtimeID
	reservationDTO := ToReservationResponse(reservation)
	result.Seats = reservationDTO.Items
	showtime := reservation.Showtime
	if showtime == nil {
		return result
	}
	result.MovieID = showtime.MovieID
	result.StartsAt = helper.TimeToPtr(showtime.StartsAt)
	result.EndsAt = helper.TimeToPtr(showtime.EndsAt)
	if showtime.Movie != nil {
		result.MovieTitle = showtime.Movie.Title
	}
	if showtime.Hall != nil {
		result.HallID = showtime.HallID
		result.HallName = showtime.Hall.Name
		result.CinemaID = showtime.Hall.CinemaID
		if showtime.Hall.Cinema != nil {
			result.CinemaName = showtime.Hall.Cinema.Name
		}
	}
	return result
}

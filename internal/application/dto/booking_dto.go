package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CreateHallRequest struct {
	Name   string          `json:"name" binding:"required,notblank,min=1,max=100"`
	Status enum.HallStatus `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type UpdateHallRequest struct {
	Name   *string          `json:"name" binding:"omitempty,notblank,min=1,max=100"`
	Status *enum.HallStatus `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type CreateSeatRequest struct {
	RowLabel string          `json:"row_label" binding:"required,notblank,max=10"`
	Number   int             `json:"number" binding:"required,gt=0,lte=1000"`
	Type     enum.SeatType   `json:"type" binding:"required,oneof=STANDARD VIP WHEELCHAIR"`
	Status   enum.SeatStatus `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type CreateSeatsRequest struct {
	Seats []CreateSeatRequest `json:"seats" binding:"required,min=1,max=500,dive"`
}

type UpdateSeatRequest struct {
	RowLabel *string          `json:"row_label" binding:"omitempty,notblank,max=10"`
	Number   *int             `json:"number" binding:"omitempty,gt=0,lte=1000"`
	Type     *enum.SeatType   `json:"type" binding:"omitempty,oneof=STANDARD VIP WHEELCHAIR"`
	Status   *enum.SeatStatus `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type SeatResponse struct {
	ID        int64           `json:"id"`
	HallID    int64           `json:"hall_id"`
	RowLabel  string          `json:"row_label"`
	Number    int             `json:"number"`
	Type      enum.SeatType   `json:"type"`
	Status    enum.SeatStatus `json:"status"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
}

type HallResponse struct {
	ID        int64           `json:"id"`
	CinemaID  int64           `json:"cinema_id"`
	Name      string          `json:"name"`
	Status    enum.HallStatus `json:"status"`
	SeatCount int             `json:"seat_count"`
	Seats     []SeatResponse  `json:"seats,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
}

type SeatTypePriceRequest struct {
	SeatType enum.SeatType `json:"seat_type" binding:"required,oneof=STANDARD VIP WHEELCHAIR"`
	Price    int64         `json:"price" binding:"required,gt=0"`
}

type CreateShowtimeRequest struct {
	MovieID    int64                  `json:"movie_id" binding:"required,gt=0"`
	StartsAt   time.Time              `json:"starts_at" binding:"required"`
	Status     enum.ShowtimeStatus    `json:"status" binding:"required,oneof=DRAFT PUBLISHED"`
	BasePrice  int64                  `json:"base_price" binding:"required,gt=0"`
	TypePrices []SeatTypePriceRequest `json:"type_prices" binding:"omitempty,max=3,dive"`
}

type UpdateShowtimeRequest struct {
	StartsAt *time.Time           `json:"starts_at"`
	Status   *enum.ShowtimeStatus `json:"status" binding:"omitempty,oneof=DRAFT PUBLISHED CANCELLED COMPLETED"`
}

type ShowtimeFilterRequest struct {
	PaginationRequest
	CinemaID *int64               `form:"cinema_id" binding:"omitempty,gt=0"`
	HallID   *int64               `form:"hall_id" binding:"omitempty,gt=0"`
	MovieID  *int64               `form:"movie_id" binding:"omitempty,gt=0"`
	From     *time.Time           `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To       *time.Time           `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
}

type ShowtimeResponse struct {
	ID         int64               `json:"id"`
	HallID     int64               `json:"hall_id"`
	HallName   string              `json:"hall_name,omitempty"`
	CinemaID   int64               `json:"cinema_id,omitempty"`
	CinemaName string              `json:"cinema_name,omitempty"`
	MovieID    int64               `json:"movie_id"`
	MovieTitle string              `json:"movie_title,omitempty"`
	StartsAt   time.Time           `json:"starts_at"`
	EndsAt     time.Time           `json:"ends_at"`
	Status     enum.ShowtimeStatus `json:"status"`
	SeatCount  int                 `json:"seat_count"`
	CreatedAt  *time.Time          `json:"created_at,omitempty"`
	UpdatedAt  *time.Time          `json:"updated_at,omitempty"`
}

type ShowtimeSeatPriceRequest struct {
	ShowtimeSeatID int64 `json:"showtime_seat_id" binding:"required,gt=0"`
	Price          int64 `json:"price" binding:"required,gt=0"`
}

type UpdateShowtimePricesRequest struct {
	Seats []ShowtimeSeatPriceRequest `json:"seats" binding:"required,min=1,max=500,dive"`
}

type ShowtimeSeatResponse struct {
	ID        int64                   `json:"id"`
	SeatID    int64                   `json:"seat_id"`
	RowLabel  string                  `json:"row_label"`
	Number    int                     `json:"number"`
	Type      enum.SeatType           `json:"type"`
	Price     int64                   `json:"price"`
	Status    enum.ShowtimeSeatStatus `json:"status"`
	HeldUntil *time.Time              `json:"held_until,omitempty"`
}

type CreateReservationRequest struct {
	ShowtimeID      int64   `json:"showtime_id" binding:"required,gt=0"`
	ShowtimeSeatIDs []int64 `json:"showtime_seat_ids" binding:"required,min=1,max=10,dive,gt=0"`
}

type ReservationItemResponse struct {
	ShowtimeSeatID int64         `json:"showtime_seat_id"`
	SeatID         int64         `json:"seat_id"`
	RowLabel       string        `json:"row_label"`
	Number         int           `json:"number"`
	Type           enum.SeatType `json:"type"`
	UnitPrice      int64         `json:"unit_price"`
}

type ReservationResponse struct {
	ID          int64                     `json:"id"`
	UserID      int64                     `json:"user_id"`
	ShowtimeID  int64                     `json:"showtime_id"`
	Status      enum.ReservationStatus    `json:"status"`
	ExpiresAt   time.Time                 `json:"expires_at"`
	TotalAmount int64                     `json:"total_amount"`
	Items       []ReservationItemResponse `json:"items"`
	CreatedAt   *time.Time                `json:"created_at,omitempty"`
	UpdatedAt   *time.Time                `json:"updated_at,omitempty"`
}

func ToSeatResponse(seat *entity.Seat) SeatResponse {
	return SeatResponse{ID: seat.ID, HallID: seat.HallID, RowLabel: seat.RowLabel, Number: seat.Number, Type: seat.Type, Status: seat.Status, CreatedAt: helper.TimeToPtr(seat.CreatedAt), UpdatedAt: helper.TimeToPtr(seat.UpdatedAt)}
}

func ToSeatResponses(seats []entity.Seat) []SeatResponse {
	result := make([]SeatResponse, 0, len(seats))
	for i := range seats {
		result = append(result, ToSeatResponse(&seats[i]))
	}
	return result
}

func ToHallResponse(hall *entity.Hall) *HallResponse {
	if hall == nil {
		return nil
	}
	return &HallResponse{ID: hall.ID, CinemaID: hall.CinemaID, Name: hall.Name, Status: hall.Status, SeatCount: len(hall.Seats), Seats: ToSeatResponses(hall.Seats), CreatedAt: helper.TimeToPtr(hall.CreatedAt), UpdatedAt: helper.TimeToPtr(hall.UpdatedAt)}
}

func ToShowtimeResponse(showtime *entity.Showtime) *ShowtimeResponse {
	if showtime == nil {
		return nil
	}
	response := &ShowtimeResponse{ID: showtime.ID, HallID: showtime.HallID, MovieID: showtime.MovieID, StartsAt: showtime.StartsAt, EndsAt: showtime.EndsAt, Status: showtime.Status, SeatCount: len(showtime.Seats), CreatedAt: helper.TimeToPtr(showtime.CreatedAt), UpdatedAt: helper.TimeToPtr(showtime.UpdatedAt)}
	if showtime.Hall != nil {
		response.HallName = showtime.Hall.Name
		response.CinemaID = showtime.Hall.CinemaID
		if showtime.Hall.Cinema != nil {
			response.CinemaName = showtime.Hall.Cinema.Name
		}
	}
	if showtime.Movie != nil {
		response.MovieTitle = showtime.Movie.Title
	}
	return response
}

func ToShowtimeSeatResponses(seats []entity.ShowtimeSeat) []ShowtimeSeatResponse {
	result := make([]ShowtimeSeatResponse, 0, len(seats))
	for i := range seats {
		item := ShowtimeSeatResponse{ID: seats[i].ID, SeatID: seats[i].SeatID, RowLabel: seats[i].RowLabel, Number: seats[i].SeatNumber, Type: seats[i].SeatType, Price: seats[i].Price, Status: seats[i].Status, HeldUntil: seats[i].HeldUntil}
		result = append(result, item)
	}
	return result
}

func ToReservationResponse(reservation *entity.Reservation) *ReservationResponse {
	if reservation == nil {
		return nil
	}
	items := make([]ReservationItemResponse, 0, len(reservation.Items))
	for i := range reservation.Items {
		item := ReservationItemResponse{ShowtimeSeatID: reservation.Items[i].ShowtimeSeatID, UnitPrice: reservation.Items[i].UnitPrice}
		if ss := reservation.Items[i].ShowtimeSeat; ss != nil {
			item.SeatID = ss.SeatID
			item.RowLabel = ss.RowLabel
			item.Number = ss.SeatNumber
			item.Type = ss.SeatType
		}
		items = append(items, item)
	}
	return &ReservationResponse{ID: reservation.ID, UserID: reservation.UserID, ShowtimeID: reservation.ShowtimeID, Status: reservation.Status, ExpiresAt: reservation.ExpiresAt, TotalAmount: reservation.TotalAmount, Items: items, CreatedAt: helper.TimeToPtr(reservation.CreatedAt), UpdatedAt: helper.TimeToPtr(reservation.UpdatedAt)}
}

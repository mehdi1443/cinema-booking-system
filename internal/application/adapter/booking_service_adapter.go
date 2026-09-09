package service_contract

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type HallServiceInterface interface {
	CreateHall(ctx context.Context, cinemaID int64, req *dto.CreateHallRequest) (*dto.HallResponse, error)
	GetHall(ctx context.Context, id int64) (*dto.HallResponse, error)
	ListHalls(ctx context.Context, cinemaID int64) ([]*dto.HallResponse, error)
	UpdateHall(ctx context.Context, id int64, req *dto.UpdateHallRequest) (*dto.HallResponse, error)
	DeleteHall(ctx context.Context, id int64) error
	CreateSeats(ctx context.Context, hallID int64, req *dto.CreateSeatsRequest) ([]dto.SeatResponse, error)
	ListSeats(ctx context.Context, hallID int64) ([]dto.SeatResponse, error)
	UpdateSeat(ctx context.Context, id int64, req *dto.UpdateSeatRequest) (*dto.SeatResponse, error)
	DeleteSeat(ctx context.Context, id int64) error
}

type ShowtimeServiceInterface interface {
	Create(ctx context.Context, hallID int64, req *dto.CreateShowtimeRequest) (*dto.ShowtimeResponse, error)
	Get(ctx context.Context, id int64) (*dto.ShowtimeResponse, error)
	List(ctx context.Context, req *dto.ShowtimeFilterRequest) ([]*dto.ShowtimeResponse, int64, error)
	Update(ctx context.Context, id int64, req *dto.UpdateShowtimeRequest) (*dto.ShowtimeResponse, error)
	Delete(ctx context.Context, id int64) error
	UpdatePrices(ctx context.Context, id int64, req *dto.UpdateShowtimePricesRequest) error
	ListSeats(ctx context.Context, id int64) ([]dto.ShowtimeSeatResponse, error)
}

type ReservationServiceInterface interface {
	Create(ctx context.Context, userID int64, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	Get(ctx context.Context, id, userID int64, allowAny bool) (*dto.ReservationResponse, error)
	ListMine(ctx context.Context, userID int64, page, limit int) ([]*dto.ReservationResponse, int64, error)
	Cancel(ctx context.Context, id, userID int64, allowAny bool) (*dto.ReservationResponse, error)
	ExpirePending(ctx context.Context, now time.Time) error
}

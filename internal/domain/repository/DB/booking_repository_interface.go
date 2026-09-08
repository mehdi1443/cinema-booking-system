package DB

import (
	"context"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type ShowtimeFilter struct {
	CinemaID   *int64
	HallID     *int64
	MovieID    *int64
	Status     *enum.ShowtimeStatus
	From       *time.Time
	To         *time.Time
	PublicOnly bool
}

type HallRepositoryInterface interface {
	CreateHall(ctx context.Context, hall *entity.Hall) (*entity.Hall, error)
	GetHall(ctx context.Context, id int64) (*entity.Hall, error)
	ListHalls(ctx context.Context, cinemaID int64) ([]entity.Hall, error)
	UpdateHall(ctx context.Context, hall *entity.Hall) (*entity.Hall, error)
	DeleteHall(ctx context.Context, id int64) error
	CreateSeats(ctx context.Context, hallID int64, seats []entity.Seat) ([]entity.Seat, error)
	GetSeat(ctx context.Context, id int64) (*entity.Seat, error)
	ListSeats(ctx context.Context, hallID int64) ([]entity.Seat, error)
	UpdateSeat(ctx context.Context, seat *entity.Seat) (*entity.Seat, error)
	DeleteSeat(ctx context.Context, id int64) error
	IsHallOwnedBy(ctx context.Context, hallID, userID int64) (bool, error)
	IsSeatOwnedBy(ctx context.Context, seatID, userID int64) (bool, error)
}

type ShowtimeRepositoryInterface interface {
	CreateShowtime(ctx context.Context, showtime *entity.Showtime, prices map[enum.SeatType]int64) (*entity.Showtime, error)
	GetShowtime(ctx context.Context, id int64) (*entity.Showtime, error)
	ListShowtimes(ctx context.Context, filter ShowtimeFilter, opts *QueryOptions) ([]entity.Showtime, int64, error)
	UpdateShowtime(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error)
	DeleteShowtime(ctx context.Context, id int64) error
	UpdateShowtimePrices(ctx context.Context, showtimeID int64, prices map[int64]int64) error
	ListShowtimeSeats(ctx context.Context, showtimeID int64) ([]entity.ShowtimeSeat, error)
	IsShowtimeOwnedBy(ctx context.Context, showtimeID, userID int64) (bool, error)
}

type ReservationRepositoryInterface interface {
	CreateReservation(ctx context.Context, userID, showtimeID int64, showtimeSeatIDs []int64, expiresAt time.Time) (*entity.Reservation, error)
	GetReservation(ctx context.Context, id int64) (*entity.Reservation, error)
	ListUserReservations(ctx context.Context, userID int64) ([]entity.Reservation, error)
	CancelReservation(ctx context.Context, id int64) (*entity.Reservation, error)
	ExpireReservations(ctx context.Context, now time.Time) error
}

type BookingRepositoryInterface interface {
	HallRepositoryInterface
	ShowtimeRepositoryInterface
	ReservationRepositoryInterface
}

package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CinemaFilter struct {
	SearchKeyword string
	ProvinceID    *int64
	CityID        *int64
	FacilityIDs   []int64
	Status        *enum.CinemaStatus
	MinRating     *float64
	MaxRating     *float64
}

type CinemaInterface interface {
	Create(ctx context.Context, cinema *entity.Cinema, facilityIDs []int64) (*entity.Cinema, error)
	GetByID(ctx context.Context, id int64) (*entity.Cinema, error)

	GetCinemas(ctx context.Context, filter CinemaFilter, opts *QueryOptions) ([]entity.Cinema, int64, error)

	Update(ctx context.Context, cinema *entity.Cinema, facilityIDs []int64) (*entity.Cinema, error)
	Delete(ctx context.Context, id int64) error

	UpdateRatingAggregate(ctx context.Context, cinemaID int64) error
}

package DB

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"golang.org/x/net/context"
)

type CinemaFilter struct {
	ProvinceID  *int64
	CityID      *int64
	FacilityIDs []int64
	Status      *enum.CinemaStatus
	Search      *string
	Limit       int
	Offset      int
}

type CinemaRepository interface {
	Create(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error)
	Update(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error)
	GetByID(ctx context.Context, id int64) (*entity.Cinema, error)
	Delete(ctx context.Context, id int64) error

	GetWithFilters(ctx context.Context, filter CinemaFilter) ([]entity.Cinema, error)

	CountWithFilters(ctx context.Context, filter CinemaFilter) (int64, error)

	AddFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	RemoveFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	UpdateFacilities(ctx context.Context, cinemaID int64, facilityIDs []int64) error
}

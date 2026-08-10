package DB

import "golang.org/x/net/context"

type CinemaFilter struct {
	ProvinceID  *int64
	CityID      *int64
	FacilityIDs []int64
	Status      *CinemaStatus
	Search      *string
	Limit       int
	Offset      int
}

type CinemaRepository interface {
	Create(ctx context.Context, cinema *Cinema) (*Cinema, error)
	Update(ctx context.Context, cinema *Cinema) (*Cinema, error)
	GetByID(ctx context.Context, id int64) (*Cinema, error)

	GetWithFilters(ctx context.Context, filter CinemaFilter) ([]Cinema, error)

	CountWithFilters(ctx context.Context, filter CinemaFilter) (int64, error)

	AddFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	RemoveFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	UpdateFacilities(ctx context.Context, cinemaID int64, facilityIDs []int64) error
}

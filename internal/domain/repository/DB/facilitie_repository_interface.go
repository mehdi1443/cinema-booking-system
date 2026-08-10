package DB

import "golang.org/x/net/context"

type FacilityRepository interface {
	Create(ctx context.Context, facility *Facility) (*Facility, error)
	GetByID(ctx context.Context, id int64) (*Facility, error)
	GetAll(ctx context.Context) ([]Facility, error)
	GetByCinemaID(ctx context.Context, cinemaID int64) ([]Facility, error)
	Update(ctx context.Context, facility *Facility) (*Facility, error)
	Delete(ctx context.Context, id int64) error
}

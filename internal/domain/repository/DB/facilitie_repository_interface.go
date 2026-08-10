package DB

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"golang.org/x/net/context"
)

type FacilityRepository interface {
	Create(ctx context.Context, facility *entity.Facility) (*entity.Facility, error)
	GetByID(ctx context.Context, id int64) (*entity.Facility, error)
	GetAll(ctx context.Context) ([]entity.Facility, error)
	GetByCinemaID(ctx context.Context, cinemaID int64) ([]entity.Facility, error)
	Update(ctx context.Context, facility *entity.Facility) (*entity.Facility, error)
	Delete(ctx context.Context, id int64) error
}

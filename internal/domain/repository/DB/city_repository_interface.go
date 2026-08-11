package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type CityRepository interface {
	Create(ctx context.Context, city *entity.City) error
	GetByID(ctx context.Context, id int64) (*entity.City, error)
	GetByNameAndProvinceID(ctx context.Context, name string, provinceID int64) (*entity.City, error)
	GetByProvinceID(ctx context.Context, provinceID int64) ([]entity.City, error)
	GetAll(ctx context.Context) ([]entity.City, error)
	Update(ctx context.Context, city *entity.City) (*entity.City, error)
	Delete(ctx context.Context, id int64) error
}

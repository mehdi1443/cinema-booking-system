package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type ProvinceRepository interface {
	Create(ctx context.Context, province *entity.Province) error
	GetByID(ctx context.Context, id int64) (*entity.Province, error)
	GetByName(ctx context.Context, name string) (*entity.Province, error)
	GetAll(ctx context.Context) ([]entity.Province, error)
	Update(ctx context.Context, province *entity.Province) (*entity.Province, error)
	Delete(ctx context.Context, id int64) error
}

package DB

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"golang.org/x/net/context"
)

type ProvinceRepository interface {
	Create(ctx context.Context, province *entity.Province) error
	GetByID(ctx context.Context, id int64) (*entity.Province, error)
	GetByName(ctx context.Context, name string) (*entity.Province, error)
	GetAll(ctx context.Context) ([]entity.Province, error)
}

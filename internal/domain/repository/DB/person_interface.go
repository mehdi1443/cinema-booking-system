package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type PersonInterface interface {
	Create(ctx context.Context, person *entity.Person) (*entity.Person, error)
	GetByID(ctx context.Context, id int64) (*entity.Person, error)
	GetByName(ctx context.Context, name string) (*entity.Person, error)
	GetAll(ctx context.Context) ([]entity.Person, error)
	Update(ctx context.Context, person *entity.Person) (*entity.Person, error)
	Delete(ctx context.Context, id int64) error
}

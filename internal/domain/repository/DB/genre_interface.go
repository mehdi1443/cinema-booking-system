package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type GenreInterface interface {
	Create(ctx context.Context, genre *entity.Genre) (*entity.Genre, error)
	GetByID(ctx context.Context, id int64) (*entity.Genre, error)
	GetByName(ctx context.Context, name string) (*entity.Genre, error)
	GetAll(ctx context.Context) ([]entity.Genre, error)
	Update(ctx context.Context, genre *entity.Genre) (*entity.Genre, error)
	Delete(ctx context.Context, id int64) error
}

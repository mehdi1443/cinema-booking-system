package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type GenreFilter struct {
	SearchKeyword string
}

type GenreInterface interface {
	Create(ctx context.Context, genre *entity.Genre) (*entity.Genre, error)
	GetByID(ctx context.Context, id int64) (*entity.Genre, error)
	GetByName(ctx context.Context, name string) (*entity.Genre, error)
	GetGenres(ctx context.Context, filter GenreFilter, opts *QueryOptions) ([]entity.Genre, int64, error)
	Update(ctx context.Context, genre *entity.Genre) (*entity.Genre, error)
	Delete(ctx context.Context, id int64) error
}

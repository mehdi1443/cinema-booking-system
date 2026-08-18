package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type MovieInterface interface {
	Create(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error)
	GetByID(ctx context.Context, id int64) (*entity.Movie, error)
	GetByStatus(ctx context.Context, status enum.MovieStatus) ([]entity.Movie, error)
	GetByGenreID(ctx context.Context, genreID int64) ([]entity.Movie, error)
	GetFiltered(ctx context.Context, filter dto.MovieFilterRequest) ([]entity.Movie, error)
	GetAll(ctx context.Context) ([]entity.Movie, error)
	Update(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error)
	Delete(ctx context.Context, id int64) error

	UpdateRatingAggregate(ctx context.Context, movieID int64) error
}

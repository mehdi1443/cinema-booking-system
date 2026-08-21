package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type MovieFilter struct {
	SearchKeyword string
	IDs           []int64
	Status        *enum.MovieStatus
	GenreID       *int64
	PersonID      *int64
	PersonRole    *enum.MoviePersonRole
	MinRating     *float64
	ReleaseYear   *int
	Country       *string
}

type MovieInterFace interface {
	Create(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error)
	GetByID(ctx context.Context, id int64) (*entity.Movie, error)

	GetMovies(ctx context.Context, filter MovieFilter, opts *QueryOptions) ([]entity.Movie, int64, error)

	Update(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error)
	Delete(ctx context.Context, id int64) error

	UpdateRatingAggregate(ctx context.Context, movieID int64) error
}

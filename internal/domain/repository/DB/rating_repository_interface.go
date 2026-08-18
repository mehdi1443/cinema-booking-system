package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type RatingInterface interface {
	Upsert(ctx context.Context, rating *entity.Rating) (*entity.Rating, error)
	GetByID(ctx context.Context, id int64) (*entity.Rating, error)
	Delete(ctx context.Context, id int64) error
	GetUserRatingForMovie(ctx context.Context, userID, movieID int64) (*entity.Rating, error)
	GetUserRatingForCinema(ctx context.Context, userID, cinemaID int64) (*entity.Rating, error)
	GetByMovieID(ctx context.Context, movieID int64, limit, offset int) ([]entity.Rating, error)
	GetByCinemaID(ctx context.Context, cinemaID int64, limit, offset int) ([]entity.Rating, error)
}

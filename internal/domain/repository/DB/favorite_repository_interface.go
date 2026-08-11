package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type FavoriteFilter struct {
	UserID int64
	Target string
	Limit  int
	Offset int
}

type FavoriteRepository interface {
	AddMovieToFavorites(ctx context.Context, userID int64, movieID int64) (*entity.Favorite, error)
	RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error
	IsMovieFavorite(ctx context.Context, userID int64, movieID int64) (bool, error)

	AddCinemaToFavorites(ctx context.Context, userID int64, cinemaID int64) (*entity.Favorite, error)
	RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error
	IsCinemaFavorite(ctx context.Context, userID int64, cinemaID int64) (bool, error)

	GetUserFavorites(ctx context.Context, filter FavoriteFilter) ([]entity.Favorite, error)
	CountUserFavorites(ctx context.Context, filter FavoriteFilter) (int64, error)
}

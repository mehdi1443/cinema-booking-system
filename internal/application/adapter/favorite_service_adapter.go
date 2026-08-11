package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type FavoriteServiceInterface interface {
	AddMovieToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteMovieRequest) (*dto.FavoriteResponse, error)
	RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error

	AddCinemaToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteCinemaRequest) (*dto.FavoriteResponse, error)
	RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error

	GetUserFavorites(ctx context.Context, userID int64, req *dto.FavoriteFilterRequest) (*dto.FavoriteListWithCountResponse, error)
}

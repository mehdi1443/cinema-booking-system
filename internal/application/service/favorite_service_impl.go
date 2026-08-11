package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type favoriteService struct {
	favoriteRepo repository_contract.FavoriteRepository
}

func NewFavoriteService(
	favoriteRepo repository_contract.FavoriteRepository,
) service_contract.FavoriteServiceInterface {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
	}
}

func (s *favoriteService) AddMovieToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteMovieRequest) (*dto.FavoriteResponse, error) {
	favorite, err := s.favoriteRepo.AddMovieToFavorites(ctx, userID, req.MovieID)
	if err != nil {
		return nil, fmt.Errorf("failed to add movie to favorites: %w", err)
	}

	return dto.ToFavoriteResponse(favorite), nil
}

func (s *favoriteService) RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error {
	if err := s.favoriteRepo.RemoveMovieFromFavorites(ctx, userID, movieID); err != nil {
		return fmt.Errorf("failed to remove movie from favorites: %w", err)
	}
	return nil
}

func (s *favoriteService) AddCinemaToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteCinemaRequest) (*dto.FavoriteResponse, error) {
	favorite, err := s.favoriteRepo.AddCinemaToFavorites(ctx, userID, req.CinemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to add cinema to favorites: %w", err)
	}

	return dto.ToFavoriteResponse(favorite), nil
}

func (s *favoriteService) RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error {
	if err := s.favoriteRepo.RemoveCinemaFromFavorites(ctx, userID, cinemaID); err != nil {
		return fmt.Errorf("failed to remove cinema from favorites: %w", err)
	}
	return nil
}

func (s *favoriteService) GetUserFavorites(ctx context.Context, userID int64, req *dto.FavoriteFilterRequest) (*dto.FavoriteListWithCountResponse, error) {
	filter := repository_contract.FavoriteFilter{
		UserID: userID,
		Target: req.Target,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	favorites, err := s.favoriteRepo.GetUserFavorites(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user favorites: %w", err)
	}

	total, err := s.favoriteRepo.CountUserFavorites(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count user favorites: %w", err)
	}

	return &dto.FavoriteListWithCountResponse{
		Favorites: dto.ToFavoriteResponseSlice(favorites),
		Total:     total,
	}, nil
}

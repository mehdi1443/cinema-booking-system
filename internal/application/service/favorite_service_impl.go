package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type favoriteService struct {
	favoriteRepo repository_contract.FavoriteRepository
	log          domainlogger.Logger
}

func NewFavoriteService(
	favoriteRepo repository_contract.FavoriteRepository,
	log domainlogger.Logger,
) service_contract.FavoriteServiceInterface {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		log:          log,
	}
}

func (s *favoriteService) AddMovieToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteMovieRequest) (*dto.FavoriteResponse, error) {
	if req == nil || userID <= 0 || req.MovieID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	favorite, err := s.favoriteRepo.AddMovieToFavorites(ctx, userID, req.MovieID)
	if err != nil {
		return nil, fmt.Errorf("add movie to favorites: %w", err)
	}

	s.log.Info(ctx, "Movie added to favorites", domainlogger.Int64("user_id", userID), domainlogger.Int64("movie_id", req.MovieID))

	return dto.ToFavoriteResponse(favorite), nil
}

func (s *favoriteService) RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error {
	if userID <= 0 || movieID <= 0 {
		return exception.ErrInvalidInput
	}

	if err := s.favoriteRepo.RemoveMovieFromFavorites(ctx, userID, movieID); err != nil {
		return fmt.Errorf("remove movie from favorites %d: %w", movieID, err)
	}

	s.log.Info(ctx, "Movie removed from favorites", domainlogger.Int64("user_id", userID), domainlogger.Int64("movie_id", movieID))

	return nil
}

func (s *favoriteService) AddCinemaToFavorites(ctx context.Context, userID int64, req *dto.AddFavoriteCinemaRequest) (*dto.FavoriteResponse, error) {
	if req == nil || userID <= 0 || req.CinemaID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	favorite, err := s.favoriteRepo.AddCinemaToFavorites(ctx, userID, req.CinemaID)
	if err != nil {
		return nil, fmt.Errorf("add cinema to favorites: %w", err)
	}

	s.log.Info(ctx, "Cinema added to favorites", domainlogger.Int64("user_id", userID), domainlogger.Int64("cinema_id", req.CinemaID))

	return dto.ToFavoriteResponse(favorite), nil
}

func (s *favoriteService) RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error {
	if userID <= 0 || cinemaID <= 0 {
		return exception.ErrInvalidInput
	}

	if err := s.favoriteRepo.RemoveCinemaFromFavorites(ctx, userID, cinemaID); err != nil {
		return fmt.Errorf("remove cinema from favorites %d: %w", cinemaID, err)
	}

	s.log.Info(ctx, "Cinema removed from favorites", domainlogger.Int64("user_id", userID), domainlogger.Int64("cinema_id", cinemaID))

	return nil
}

func (s *favoriteService) GetUserFavorites(ctx context.Context, userID int64, req *dto.FavoriteFilterRequest) (*dto.FavoriteListWithCountResponse, error) {
	if req == nil || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	filter := repository_contract.FavoriteFilter{
		UserID: userID,
		Target: req.Target,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	favorites, err := s.favoriteRepo.GetUserFavorites(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get user favorites: %w", err)
	}

	total, err := s.favoriteRepo.CountUserFavorites(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("count user favorites: %w", err)
	}

	return &dto.FavoriteListWithCountResponse{
		Favorites: dto.ToFavoriteResponseSlice(favorites),
		Total:     total,
	}, nil
}

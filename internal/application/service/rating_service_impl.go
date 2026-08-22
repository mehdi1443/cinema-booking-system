package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type ratingService struct {
	ratingRepo repository_contract.RatingInterface
	movieRepo  repository_contract.MovieInterFace
	cinemaRepo repository_contract.CinemaInterface
}

func NewRatingService(
	ratingRepo repository_contract.RatingInterface,
	movieRepo repository_contract.MovieInterFace,
	cinemaRepo repository_contract.CinemaInterface,
) service_contract.RatingServiceInterface {
	return &ratingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
		cinemaRepo: cinemaRepo,
	}
}

func (s *ratingService) SubmitRating(ctx context.Context, userID int64, req *dto.SubmitRatingRequest) (*dto.RatingResponse, error) {
	if (req.MovieID == nil && req.CinemaID == nil) || (req.MovieID != nil && req.CinemaID != nil) {
		return nil, exception.ErrInvalidInput
	}

	rating := &entity.Rating{
		UserID:   userID,
		MovieID:  req.MovieID,
		CinemaID: req.CinemaID,
		Score:    req.Score,
	}

	submittedRating, err := s.ratingRepo.Upsert(ctx, rating)
	if err != nil {
		return nil, err
	}

	if req.MovieID != nil {
		if err := s.movieRepo.UpdateRatingAggregate(ctx, *req.MovieID); err != nil {
			return nil, err
		}
	} else if req.CinemaID != nil {
		if err := s.cinemaRepo.UpdateRatingAggregate(ctx, *req.CinemaID); err != nil {
			return nil, err
		}
	}

	return dto.ToRatingResponse(submittedRating), nil
}

func (s *ratingService) GetUserRatingForMovie(ctx context.Context, userID, movieID int64) (*dto.RatingResponse, error) {
	rating, err := s.ratingRepo.GetUserRatingForMovie(ctx, userID, movieID)
	if err != nil {
		return nil, err
	}

	return dto.ToRatingResponse(rating), nil
}

func (s *ratingService) GetUserRatingForCinema(ctx context.Context, userID, cinemaID int64) (*dto.RatingResponse, error) {
	rating, err := s.ratingRepo.GetUserRatingForCinema(ctx, userID, cinemaID)
	if err != nil {
		return nil, err
	}

	return dto.ToRatingResponse(rating), nil
}

func (s *ratingService) GetByMovieID(ctx context.Context, movieID int64, page, limit int) ([]*dto.RatingResponse, error) {
	_, limit, offset := normalizePagination(page, limit)

	ratings, err := s.ratingRepo.GetByMovieID(ctx, movieID, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.ToRatingResponseSlice(ratings), nil
}

func (s *ratingService) GetByCinemaID(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.RatingResponse, error) {
	_, limit, offset := normalizePagination(page, limit)

	ratings, err := s.ratingRepo.GetByCinemaID(ctx, cinemaID, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.ToRatingResponseSlice(ratings), nil
}

func (s *ratingService) Delete(ctx context.Context, userID int64, ratingID int64) error {
	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return err
	}

	if rating.UserID != userID {
		return exception.ErrUnauthorized
	}

	if err := s.ratingRepo.Delete(ctx, ratingID); err != nil {
		return err
	}

	if rating.MovieID != nil {
		return s.movieRepo.UpdateRatingAggregate(ctx, *rating.MovieID)
	} else if rating.CinemaID != nil {
		return s.cinemaRepo.UpdateRatingAggregate(ctx, *rating.CinemaID)
	}

	return nil
}

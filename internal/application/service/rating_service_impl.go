package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type ratingService struct {
	ratingRepo repository_contract.RatingInterface
	movieRepo  repository_contract.MovieInterFace
	cinemaRepo repository_contract.CinemaInterface
	log        domainlogger.Logger
}

func NewRatingService(
	ratingRepo repository_contract.RatingInterface,
	movieRepo repository_contract.MovieInterFace,
	cinemaRepo repository_contract.CinemaInterface,
	log domainlogger.Logger,
) service_contract.RatingServiceInterface {
	return &ratingService{
		ratingRepo: ratingRepo,
		movieRepo:  movieRepo,
		cinemaRepo: cinemaRepo,
		log:        log,
	}
}

func (s *ratingService) SubmitRating(ctx context.Context, userID int64, req *dto.SubmitRatingRequest) (*dto.RatingResponse, error) {
	if req == nil || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}

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
		return nil, fmt.Errorf("submit rating: %w", err)
	}

	if req.MovieID != nil {
		if err := s.movieRepo.UpdateRatingAggregate(ctx, *req.MovieID); err != nil {
			return nil, fmt.Errorf("update movie rating aggregate: %w", err)
		}
		s.log.Info(ctx, "Movie rating submitted", domainlogger.Int64("user_id", userID), domainlogger.Int64("movie_id", *req.MovieID), domainlogger.Int("score", int(req.Score)))
	} else if req.CinemaID != nil {
		if err := s.cinemaRepo.UpdateRatingAggregate(ctx, *req.CinemaID); err != nil {
			return nil, fmt.Errorf("update cinema rating aggregate: %w", err)
		}
		s.log.Info(ctx, "Cinema rating submitted", domainlogger.Int64("user_id", userID), domainlogger.Int64("cinema_id", *req.CinemaID), domainlogger.Int("score", int(req.Score)))
	}

	return dto.ToRatingResponse(submittedRating), nil
}

func (s *ratingService) GetUserRatingForMovie(ctx context.Context, userID, movieID int64) (*dto.RatingResponse, error) {
	if userID <= 0 || movieID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	rating, err := s.ratingRepo.GetUserRatingForMovie(ctx, userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("get user rating for movie %d: %w", movieID, err)
	}

	return dto.ToRatingResponse(rating), nil
}

func (s *ratingService) GetUserRatingForCinema(ctx context.Context, userID, cinemaID int64) (*dto.RatingResponse, error) {
	if userID <= 0 || cinemaID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	rating, err := s.ratingRepo.GetUserRatingForCinema(ctx, userID, cinemaID)
	if err != nil {
		return nil, fmt.Errorf("get user rating for cinema %d: %w", cinemaID, err)
	}

	return dto.ToRatingResponse(rating), nil
}

func (s *ratingService) GetByMovieID(ctx context.Context, movieID int64, page, limit int) ([]*dto.RatingResponse, error) {
	if movieID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	_, limit, offset := normalizePagination(page, limit)

	ratings, err := s.ratingRepo.GetByMovieID(ctx, movieID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get ratings by movie %d: %w", movieID, err)
	}

	return dto.ToRatingResponseSlice(ratings), nil
}

func (s *ratingService) GetByCinemaID(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.RatingResponse, error) {
	if cinemaID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	_, limit, offset := normalizePagination(page, limit)

	ratings, err := s.ratingRepo.GetByCinemaID(ctx, cinemaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get ratings by cinema %d: %w", cinemaID, err)
	}

	return dto.ToRatingResponseSlice(ratings), nil
}

func (s *ratingService) Delete(ctx context.Context, userID int64, ratingID int64) error {
	if userID <= 0 || ratingID <= 0 {
		return exception.ErrInvalidInput
	}

	rating, err := s.ratingRepo.GetByID(ctx, ratingID)
	if err != nil {
		return fmt.Errorf("get rating before delete %d: %w", ratingID, err)
	}

	if rating.UserID != userID {
		return exception.ErrNotAccessible
	}

	if err := s.ratingRepo.Delete(ctx, ratingID); err != nil {
		return fmt.Errorf("delete rating %d: %w", ratingID, err)
	}

	if rating.MovieID != nil {
		s.log.Info(ctx, "Movie rating deleted", domainlogger.Int64("user_id", userID), domainlogger.Int64("rating_id", ratingID))
		return s.movieRepo.UpdateRatingAggregate(ctx, *rating.MovieID)
	} else if rating.CinemaID != nil {
		s.log.Info(ctx, "Cinema rating deleted", domainlogger.Int64("user_id", userID), domainlogger.Int64("rating_id", ratingID))
		return s.cinemaRepo.UpdateRatingAggregate(ctx, *rating.CinemaID)
	}

	return nil
}

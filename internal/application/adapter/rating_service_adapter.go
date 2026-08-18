package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type RatingServiceInterface interface {
	SubmitRating(ctx context.Context, userID int64, req *dto.SubmitRatingRequest) (*dto.RatingResponse, error)
	GetUserRatingForMovie(ctx context.Context, userID, movieID int64) (*dto.RatingResponse, error)
	GetUserRatingForCinema(ctx context.Context, userID, cinemaID int64) (*dto.RatingResponse, error)
	GetByMovieID(ctx context.Context, movieID int64, page, limit int) ([]*dto.RatingResponse, error)
	GetByCinemaID(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.RatingResponse, error)
	Delete(ctx context.Context, userID int64, ratingID int64) error
}

package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type MovieServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateMovieRequest) (*dto.MovieDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.MovieDetailResponse, error)
	GetMovies(ctx context.Context, req dto.MovieFilterRequest) ([]*dto.MovieListResponse, int64, error)
	Update(ctx context.Context, id int64, req *dto.UpdateMovieRequest) (*dto.MovieDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadPoster(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error)
	DeletePoster(ctx context.Context, id int64) error
	UploadTrailer(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error)
	DeleteTrailer(ctx context.Context, id int64) error
}

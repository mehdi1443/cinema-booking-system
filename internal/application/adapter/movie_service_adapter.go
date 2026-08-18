package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type MovieServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateMovieRequest) (*dto.MovieDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.MovieDetailResponse, error)
	GetByStatus(ctx context.Context, status enum.MovieStatus) ([]*dto.MovieListResponse, error)
	GetByGenreID(ctx context.Context, genreID int64) ([]*dto.MovieListResponse, error)
	GetFilteredMovies(ctx context.Context, req dto.MovieFilterRequest) ([]*dto.MovieListResponse, error)
	GetAll(ctx context.Context) ([]*dto.MovieListResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdateMovieRequest) (*dto.MovieDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadPoster(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error)
	DeletePoster(ctx context.Context, id int64) error
	UploadTrailer(ctx context.Context, id int64, file io.Reader) (*dto.MovieDetailResponse, error)
	DeleteTrailer(ctx context.Context, id int64) error
}

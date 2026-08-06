package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/dto"
)

type GenreUseCaseInterface interface {
	CreateGenre(ctx context.Context, req *dto.CreateGenreRequest) (*dto.GenreResponse, error)
	GetGenreByID(ctx context.Context, id int64) (*dto.GenreResponse, error)
	GetAllGenres(ctx context.Context) ([]*dto.GenreResponse, error)
	UpdateGenre(ctx context.Context, id int64, req *dto.UpdateGenreRequest) (*dto.GenreResponse, error)
	DeleteGenre(ctx context.Context, id int64) error
}

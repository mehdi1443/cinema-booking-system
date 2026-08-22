package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type CinemaServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateCinemaRequest) (*dto.CinemaDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.CinemaDetailResponse, error)
	GetCinemas(ctx context.Context, req dto.CinemaFilterRequest) ([]*dto.CinemaListResponse, int64, error)
	Update(ctx context.Context, id int64, req *dto.UpdateCinemaRequest) (*dto.CinemaDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadImage(ctx context.Context, id int64, file io.Reader, filename string) (*dto.CinemaDetailResponse, error)
	DeleteImage(ctx context.Context, id int64) error
}

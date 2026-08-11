package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type CinemaServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateCinemaRequest) (*dto.CinemaDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.CinemaDetailResponse, error)
	GetWithFilters(ctx context.Context, req *dto.CinemaFilterRequest) (*dto.CinemaListWithCountResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdateCinemaRequest) (*dto.CinemaDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadImage(ctx context.Context, id int64, file io.Reader, filename string) (*dto.CinemaDetailResponse, error)
	DeleteImage(ctx context.Context, id int64) error

	AddFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	RemoveFacility(ctx context.Context, cinemaID int64, facilityID int64) error
	UpdateFacilities(ctx context.Context, cinemaID int64, facilityIDs []int64) error
}

package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/dto"
)

type PersonServiceInterface interface {
	Create(ctx context.Context, req *dto.CreatePersonRequest) (*dto.PersonDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.PersonDetailResponse, error)
	GetAll(ctx context.Context) ([]*dto.PersonListResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdatePersonRequest) (*dto.PersonDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadImage(ctx context.Context, id int64, file io.Reader) (*dto.PersonDetailResponse, error)
	DeleteImage(ctx context.Context, id int64) error
}

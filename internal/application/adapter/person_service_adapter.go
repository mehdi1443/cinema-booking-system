package service_contract

import (
	"context"
	"io"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type PersonServiceInterface interface {
	Create(ctx context.Context, req *dto.CreatePersonRequest) (*dto.PersonDetailResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.PersonDetailResponse, error)
	GetPersons(ctx context.Context, req dto.PersonFilterRequest) ([]*dto.PersonListResponse, int64, error)
	Update(ctx context.Context, id int64, req *dto.UpdatePersonRequest) (*dto.PersonDetailResponse, error)
	Delete(ctx context.Context, id int64) error

	UploadImage(ctx context.Context, id int64, file io.Reader) (*dto.PersonDetailResponse, error)
	DeleteImage(ctx context.Context, id int64) error
}

package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type FacilityServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateFacilityRequest) (*dto.FacilityResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.FacilityResponse, error)
	GetAll(ctx context.Context) ([]dto.FacilityResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdateFacilityRequest) (*dto.FacilityResponse, error)
	Delete(ctx context.Context, id int64) error
}

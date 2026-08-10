package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type facilityService struct {
	facilityRepo repository_contract.FacilityRepository
}

func NewFacilityService(facilityRepo repository_contract.FacilityRepository) service_contract.FacilityServiceInterface {
	return &facilityService{
		facilityRepo: facilityRepo,
	}
}

func (s *facilityService) Create(ctx context.Context, req *dto.CreateFacilityRequest) (*dto.FacilityResponse, error) {
	facility := &entity.Facility{
		Name: req.Name,
	}

	createdFacility, err := s.facilityRepo.Create(ctx, facility)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %w", err)
	}

	return dto.ToFacilityResponse(createdFacility), nil
}

func (s *facilityService) GetByID(ctx context.Context, id int64) (*dto.FacilityResponse, error) {
	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToFacilityResponse(facility), nil
}

func (s *facilityService) GetAll(ctx context.Context) ([]dto.FacilityResponse, error) {
	facilities, err := s.facilityRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch facilities: %w", err)
	}

	return dto.ToFacilityResponseSlice(facilities), nil
}

func (s *facilityService) Update(ctx context.Context, id int64, req *dto.UpdateFacilityRequest) (*dto.FacilityResponse, error) {
	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	facility.Name = req.Name

	updatedFacility, err := s.facilityRepo.Update(ctx, facility)
	if err != nil {
		return nil, fmt.Errorf("failed to update facility: %w", err)
	}

	return dto.ToFacilityResponse(updatedFacility), nil
}

func (s *facilityService) Delete(ctx context.Context, id int64) error {
	_, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.facilityRepo.Delete(ctx, id)
}

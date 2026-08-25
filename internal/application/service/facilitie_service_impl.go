package service

import (
	"context"
	"fmt"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
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
	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	facility := &entity.Facility{
		Name: trimmedName,
	}

	createdFacility, err := s.facilityRepo.Create(ctx, facility)
	if err != nil {
		return nil, fmt.Errorf("create facility: %w", err)
	}

	return dto.ToFacilityResponse(createdFacility), nil
}

func (s *facilityService) GetByID(ctx context.Context, id int64) (*dto.FacilityResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get facility by id %d: %w", id, err)
	}

	return dto.ToFacilityResponse(facility), nil
}

func (s *facilityService) GetAll(ctx context.Context) ([]dto.FacilityResponse, error) {
	facilities, err := s.facilityRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all facilities: %w", err)
	}

	return dto.ToFacilityResponseSlice(facilities), nil
}

func (s *facilityService) Update(ctx context.Context, id int64, req *dto.UpdateFacilityRequest) (*dto.FacilityResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get facility before update %d: %w", id, err)
	}

	facility.Name = trimmedName

	updatedFacility, err := s.facilityRepo.Update(ctx, facility)
	if err != nil {
		return nil, fmt.Errorf("update facility %d: %w", id, err)
	}

	return dto.ToFacilityResponse(updatedFacility), nil
}

func (s *facilityService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	if _, err := s.facilityRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get facility before delete %d: %w", id, err)
	}

	if err := s.facilityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete facility %d: %w", id, err)
	}

	return nil
}

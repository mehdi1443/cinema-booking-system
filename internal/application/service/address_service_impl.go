package service

import (
	"context"
	"fmt"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type addressService struct {
	provinceRepo repository_contract.ProvinceRepository
	cityRepo     repository_contract.CityRepository
}

func NewAddressService(
	provinceRepo repository_contract.ProvinceRepository,
	cityRepo repository_contract.CityRepository,
) service_contract.AddressServiceInterface {
	return &addressService{
		provinceRepo: provinceRepo,
		cityRepo:     cityRepo,
	}
}

// ==================== Province Methods ====================

func (s *addressService) CreateProvince(ctx context.Context, req *dto.CreateProvinceRequest) (*dto.ProvinceResponse, error) {
	province := &entity.Province{
		Name: req.Name,
	}

	if err := s.provinceRepo.Create(ctx, province); err != nil {
		return nil, fmt.Errorf("failed to create province: %w", err)
	}

	return dto.ToProvinceResponse(province), nil
}

func (s *addressService) GetProvinceByID(ctx context.Context, id int64) (*dto.ProvinceResponse, error) {
	province, err := s.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToProvinceResponse(province), nil
}

func (s *addressService) GetAllProvinces(ctx context.Context) ([]*dto.ProvinceResponse, error) {
	provinces, err := s.provinceRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %w", err)
	}

	return dto.ToProvinceResponseSlice(provinces), nil
}

func (s *addressService) UpdateProvince(ctx context.Context, id int64, req *dto.UpdateProvinceRequest) (*dto.ProvinceResponse, error) {
	province, err := s.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	province.Name = req.Name

	updatedProvince, err := s.provinceRepo.Update(ctx, province)
	if err != nil {
		return nil, fmt.Errorf("failed to update province: %w", err)
	}

	return dto.ToProvinceResponse(updatedProvince), nil
}

func (s *addressService) DeleteProvince(ctx context.Context, id int64) error {
	if _, err := s.provinceRepo.GetByID(ctx, id); err != nil {
		return err
	}

	if err := s.provinceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete province: %w", err)
	}

	return nil
}

// ==================== City Methods ====================

func (s *addressService) CreateCity(ctx context.Context, req *dto.CreateCityRequest) (*dto.CityResponse, error) {
	city := &entity.City{
		ProvinceID: req.ProvinceID,
		Name:       req.Name,
	}

	if err := s.cityRepo.Create(ctx, city); err != nil {
		return nil, fmt.Errorf("failed to create city: %w", err)
	}

	return dto.ToCityResponse(city), nil
}

func (s *addressService) GetCityByID(ctx context.Context, id int64) (*dto.CityResponse, error) {
	city, err := s.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToCityResponse(city), nil
}

func (s *addressService) GetCitiesByProvinceID(ctx context.Context, provinceID int64) ([]*dto.CityResponse, error) {
	cities, err := s.cityRepo.GetByProvinceID(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities for province %d: %w", provinceID, err)
	}

	return dto.ToCityResponseSlice(cities), nil
}

func (s *addressService) GetAllCities(ctx context.Context) ([]*dto.CityResponse, error) {
	cities, err := s.cityRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}

	return dto.ToCityResponseSlice(cities), nil
}

func (s *addressService) UpdateCity(ctx context.Context, id int64, req *dto.UpdateCityRequest) (*dto.CityResponse, error) {
	city, err := s.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ProvinceID != nil {
		city.ProvinceID = *req.ProvinceID
	}
	if req.Name != nil {
		city.Name = *req.Name
	}

	updatedCity, err := s.cityRepo.Update(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("failed to update city: %w", err)
	}

	return dto.ToCityResponse(updatedCity), nil
}

func (s *addressService) DeleteCity(ctx context.Context, id int64) error {
	if _, err := s.cityRepo.GetByID(ctx, id); err != nil {
		return err
	}

	if err := s.cityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete city: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"fmt"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type addressService struct {
	provinceRepo repository_contract.ProvinceRepository
	cityRepo     repository_contract.CityRepository
	log          domainlogger.Logger
}

func NewAddressService(
	provinceRepo repository_contract.ProvinceRepository,
	cityRepo repository_contract.CityRepository,
	log domainlogger.Logger,
) service_contract.AddressServiceInterface {
	return &addressService{
		provinceRepo: provinceRepo,
		cityRepo:     cityRepo,
		log:          log,
	}
}

func (s *addressService) CreateProvince(ctx context.Context, req *dto.CreateProvinceRequest) (*dto.ProvinceResponse, error) {
	if req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	province := &entity.Province{
		Name: trimmedName,
	}

	if err := s.provinceRepo.Create(ctx, province); err != nil {
		return nil, fmt.Errorf("create province: %w", err)
	}

	s.log.Info(ctx, "Province created successfully", domainlogger.Int64("province_id", province.ID))

	return dto.ToProvinceResponse(province), nil
}

func (s *addressService) GetProvinceByID(ctx context.Context, id int64) (*dto.ProvinceResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	province, err := s.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get province by id %d: %w", id, err)
	}

	return dto.ToProvinceResponse(province), nil
}

func (s *addressService) GetAllProvinces(ctx context.Context) ([]*dto.ProvinceResponse, error) {
	provinces, err := s.provinceRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all provinces: %w", err)
	}

	return dto.ToProvinceResponseSlice(provinces), nil
}

func (s *addressService) UpdateProvince(ctx context.Context, id int64, req *dto.UpdateProvinceRequest) (*dto.ProvinceResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	province, err := s.provinceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get province before update %d: %w", id, err)
	}

	province.Name = trimmedName

	updatedProvince, err := s.provinceRepo.Update(ctx, province)
	if err != nil {
		return nil, fmt.Errorf("update province %d: %w", id, err)
	}

	s.log.Info(ctx, "Province updated successfully", domainlogger.Int64("province_id", updatedProvince.ID))

	return dto.ToProvinceResponse(updatedProvince), nil
}

func (s *addressService) DeleteProvince(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	if _, err := s.provinceRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get province before delete %d: %w", id, err)
	}

	if err := s.provinceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete province %d: %w", id, err)
	}

	s.log.Info(ctx, "Province deleted successfully", domainlogger.Int64("province_id", id))

	return nil
}

func (s *addressService) CreateCity(ctx context.Context, req *dto.CreateCityRequest) (*dto.CityResponse, error) {
	if req == nil || req.ProvinceID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, exception.ErrInvalidInput
	}

	city := &entity.City{
		ProvinceID: req.ProvinceID,
		Name:       trimmedName,
	}

	if err := s.cityRepo.Create(ctx, city); err != nil {
		return nil, fmt.Errorf("create city: %w", err)
	}

	s.log.Info(ctx, "City created successfully", domainlogger.Int64("city_id", city.ID))

	return dto.ToCityResponse(city), nil
}

func (s *addressService) GetCityByID(ctx context.Context, id int64) (*dto.CityResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	city, err := s.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get city by id %d: %w", id, err)
	}

	return dto.ToCityResponse(city), nil
}

func (s *addressService) GetCitiesByProvinceID(ctx context.Context, provinceID int64) ([]*dto.CityResponse, error) {
	if provinceID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	cities, err := s.cityRepo.GetByProvinceID(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("get cities for province %d: %w", provinceID, err)
	}

	return dto.ToCityResponseSlice(cities), nil
}

func (s *addressService) GetAllCities(ctx context.Context) ([]*dto.CityResponse, error) {
	cities, err := s.cityRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all cities: %w", err)
	}

	return dto.ToCityResponseSlice(cities), nil
}

func (s *addressService) UpdateCity(ctx context.Context, id int64, req *dto.UpdateCityRequest) (*dto.CityResponse, error) {
	if id <= 0 || req == nil {
		return nil, exception.ErrInvalidInput
	}

	city, err := s.cityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get city before update %d: %w", id, err)
	}

	if req.ProvinceID != nil {
		if *req.ProvinceID <= 0 {
			return nil, exception.ErrInvalidInput
		}
		city.ProvinceID = *req.ProvinceID
	}

	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, exception.ErrInvalidInput
		}
		city.Name = trimmedName
	}

	updatedCity, err := s.cityRepo.Update(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("update city %d: %w", id, err)
	}

	s.log.Info(ctx, "City updated successfully", domainlogger.Int64("city_id", updatedCity.ID))

	return dto.ToCityResponse(updatedCity), nil
}

func (s *addressService) DeleteCity(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.ErrInvalidInput
	}

	if _, err := s.cityRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get city before delete %d: %w", id, err)
	}

	if err := s.cityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete city %d: %w", id, err)
	}

	s.log.Info(ctx, "City deleted successfully", domainlogger.Int64("city_id", id))

	return nil
}

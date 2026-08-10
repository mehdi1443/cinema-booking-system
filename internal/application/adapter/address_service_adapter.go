package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
)

type AddressServiceInterface interface {
	// Province Methods
	CreateProvince(ctx context.Context, req *dto.CreateProvinceRequest) (*dto.ProvinceResponse, error)
	GetProvinceByID(ctx context.Context, id int64) (*dto.ProvinceResponse, error)
	GetAllProvinces(ctx context.Context) ([]*dto.ProvinceResponse, error)
	UpdateProvince(ctx context.Context, id int64, req *dto.UpdateProvinceRequest) (*dto.ProvinceResponse, error)
	DeleteProvince(ctx context.Context, id int64) error

	// City Methods
	CreateCity(ctx context.Context, req *dto.CreateCityRequest) (*dto.CityResponse, error)
	GetCityByID(ctx context.Context, id int64) (*dto.CityResponse, error)
	GetCitiesByProvinceID(ctx context.Context, provinceID int64) ([]*dto.CityResponse, error)
	GetAllCities(ctx context.Context) ([]*dto.CityResponse, error)
	UpdateCity(ctx context.Context, id int64, req *dto.UpdateCityRequest) (*dto.CityResponse, error)
	DeleteCity(ctx context.Context, id int64) error
}

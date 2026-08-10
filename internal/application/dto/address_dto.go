package dto

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type CreateProvinceRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateProvinceRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type CreateCityRequest struct {
	ProvinceID int64  `json:"province_id" binding:"required"`
	Name       string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateCityRequest struct {
	ProvinceID *int64  `json:"province_id"`
	Name       *string `json:"name" binding:"omitempty,min=2,max=100"`
}

type ProvinceResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	ProvinceID int64             `json:"province_id"`
	Province   *ProvinceResponse `json:"province,omitempty"`
}

func ToProvinceResponse(province *entity.Province) *ProvinceResponse {
	if province == nil {
		return nil
	}
	return &ProvinceResponse{
		ID:   province.ID,
		Name: province.Name,
	}
}

func ToProvinceResponseSlice(provinces []entity.Province) []*ProvinceResponse {
	if len(provinces) == 0 {
		return []*ProvinceResponse{}
	}

	responses := make([]*ProvinceResponse, 0, len(provinces))
	for ind := range provinces {
		if resp := ToProvinceResponse(&provinces[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

func ToCityResponse(city *entity.City) *CityResponse {
	if city == nil {
		return nil
	}

	var provinceResp *ProvinceResponse
	if city.Province != nil {
		provinceResp = ToProvinceResponse(city.Province)
	}

	return &CityResponse{
		ID:         city.ID,
		Name:       city.Name,
		ProvinceID: city.ProvinceID,
		Province:   provinceResp,
	}
}

func ToCityResponseSlice(cities []entity.City) []*CityResponse {
	if len(cities) == 0 {
		return []*CityResponse{}
	}

	responses := make([]*CityResponse, 0, len(cities))
	for ind := range cities {
		if resp := ToCityResponse(&cities[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

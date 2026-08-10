package dto

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

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

package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CinemaFilterRequest struct {
	ProvinceID  *int64             `form:"province_id"`
	CityID      *int64             `form:"city_id"`
	FacilityIDs []int64            `form:"facility_ids"`
	Status      *enum.CinemaStatus `form:"status"`
	Search      *string            `form:"search"`
	Limit       int                `form:"limit,default=10" binding:"omitempty,gte=1,lte=100"`
	Offset      int                `form:"offset,default=0" binding:"omitempty,gte=0"`
}

type CreateCinemaRequest struct {
	CityID      int64             `json:"city_id" binding:"required"`
	Name        string            `json:"name" binding:"required,min=2,max=255"`
	Description string            `json:"description"`
	Address     string            `json:"address" binding:"required"`
	Latitude    *float64          `json:"latitude"`
	Longitude   *float64          `json:"longitude"`
	Phone       string            `json:"phone"`
	Status      enum.CinemaStatus `json:"status" binding:"required"`
	FacilityIDs []int64           `json:"facility_ids"`
}

type UpdateCinemaRequest struct {
	CityID      *int64             `json:"city_id"`
	Name        *string            `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string            `json:"description"`
	Address     *string            `json:"address"`
	Latitude    *float64           `json:"latitude"`
	Longitude   *float64           `json:"longitude"`
	Phone       *string            `json:"phone"`
	Status      *enum.CinemaStatus `json:"status"`
	FacilityIDs *[]int64           `json:"facility_ids"`
}

type CinemaListResponse struct {
	ID           int64             `json:"id"`
	Name         string            `json:"name"`
	Address      string            `json:"address"`
	ImageURLLow  string            `json:"image_url_low"`
	Status       enum.CinemaStatus `json:"status"`
	CityName     string            `json:"city_name,omitempty"`
	ProvinceName string            `json:"province_name,omitempty"`
}

type CinemaDetailResponse struct {
	ID           int64              `json:"id"`
	CityID       int64              `json:"city_id"`
	City         *CityResponse      `json:"city,omitempty"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Address      string             `json:"address"`
	Latitude     *float64           `json:"latitude,omitempty"`
	Longitude    *float64           `json:"longitude,omitempty"`
	Phone        string             `json:"phone"`
	ImageURLLow  string             `json:"image_url_low"`
	ImageURLHigh string             `json:"image_url_high"`
	Status       enum.CinemaStatus  `json:"status"`
	Facilities   []FacilityResponse `json:"facilities,omitempty"`
	CreatedAt    *time.Time         `json:"created_at,omitempty"`
	UpdatedAt    *time.Time         `json:"updated_at,omitempty"`
}

type CinemaListWithCountResponse struct {
	Cinemas []CinemaListResponse `json:"cinemas"`
	Total   int64                `json:"total"`
}

func ToCinemaListResponse(cinema *entity.Cinema) *CinemaListResponse {
	if cinema == nil {
		return nil
	}

	cityName := ""
	provinceName := ""
	if cinema.City != nil {
		cityName = cinema.City.Name
		if cinema.City.Province != nil {
			provinceName = cinema.City.Province.Name
		}
	}

	var imageURLLow string
	if cinema.ImageURLLow != nil {
		imageURLLow = *cinema.ImageURLLow
	}

	return &CinemaListResponse{
		ID:           cinema.ID,
		Name:         cinema.Name,
		Address:      cinema.Address,
		ImageURLLow:  imageURLLow,
		Status:       cinema.Status,
		CityName:     cityName,
		ProvinceName: provinceName,
	}
}

func ToCinemaListResponseSlice(cinemas []entity.Cinema) []CinemaListResponse {
	if len(cinemas) == 0 {
		return []CinemaListResponse{}
	}

	responses := make([]CinemaListResponse, 0, len(cinemas))
	for ind := range cinemas {
		if resp := ToCinemaListResponse(&cinemas[ind]); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}

func ToCinemaDetailResponse(cinema *entity.Cinema) *CinemaDetailResponse {
	if cinema == nil {
		return nil
	}

	var cityResp *CityResponse
	if cinema.City != nil {
		cityResp = ToCityResponse(cinema.City)
	}

	var facilities []FacilityResponse
	if len(cinema.Facilities) > 0 {
		facilities = make([]FacilityResponse, 0, len(cinema.Facilities))
		for _, f := range cinema.Facilities {
			if resp := ToFacilityResponse(&f); resp != nil {
				facilities = append(facilities, *resp)
			}
		}
	}

	var description, phone, imageURLLow, imageURLHigh string
	if cinema.Description != nil {
		description = *cinema.Description
	}
	if cinema.Phone != nil {
		phone = *cinema.Phone
	}
	if cinema.ImageURLLow != nil {
		imageURLLow = *cinema.ImageURLLow
	}
	if cinema.ImageURLHigh != nil {
		imageURLHigh = *cinema.ImageURLHigh
	}

	return &CinemaDetailResponse{
		ID:           cinema.ID,
		CityID:       cinema.CityID,
		City:         cityResp,
		Name:         cinema.Name,
		Description:  description,
		Address:      cinema.Address,
		Latitude:     cinema.Latitude,
		Longitude:    cinema.Longitude,
		Phone:        phone,
		ImageURLLow:  imageURLLow,
		ImageURLHigh: imageURLHigh,
		Status:       cinema.Status,
		Facilities:   facilities,
		CreatedAt:    helper.TimeToPtr(cinema.CreatedAt),
		UpdatedAt:    helper.TimeToPtr(cinema.UpdatedAt),
	}
}

package dto

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type CreateFacilityRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateFacilityRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type FacilityResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToFacilityResponse(facility *entity.Facility) *FacilityResponse {
	if facility == nil {
		return nil
	}
	return &FacilityResponse{
		ID:   facility.ID,
		Name: facility.Name,
	}
}

func ToFacilityResponseSlice(facilities []entity.Facility) []FacilityResponse {
	if len(facilities) == 0 {
		return []FacilityResponse{}
	}

	responses := make([]FacilityResponse, 0, len(facilities))
	for ind := range facilities {
		if resp := ToFacilityResponse(&facilities[ind]); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}

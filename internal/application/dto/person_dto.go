package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type CreatePersonRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100"`
	Biography string `json:"biography"`
}

type UpdatePersonRequest struct {
	Name      *string `json:"name" binding:"required,min=2,max=100"`
	Biography *string `json:"biography"`
}

type PersonListResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ImageURLLow string `json:"image_url_low"`
}

type PersonDetailResponse struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Biography    string     `json:"biography"`
	ImageURLLow  string     `json:"image_url_low"`
	ImageURLHigh string     `json:"image_url_high"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func ToPersonListResponse(person *entity.Person) *PersonListResponse {
	if person == nil {
		return nil
	}

	return &PersonListResponse{
		ID:          person.ID,
		Name:        person.Name,
		ImageURLLow: person.ImageURLLow,
	}
}

func ToPersonListResponseSlice(persons []entity.Person) []*PersonListResponse {
	if len(persons) == 0 {
		return []*PersonListResponse{}
	}

	responses := make([]*PersonListResponse, 0, len(persons))
	for ind := range persons {
		if resp := ToPersonListResponse(&persons[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

func ToPersonDetailResponse(person *entity.Person) *PersonDetailResponse {
	if person == nil {
		return nil
	}

	return &PersonDetailResponse{
		ID:           person.ID,
		Name:         person.Name,
		Biography:    person.Biography,
		ImageURLLow:  person.ImageURLLow,
		ImageURLHigh: person.ImageURLHigh,
		CreatedAt:    helper.TimeToPtr(person.CreatedAt),
		UpdatedAt:    helper.TimeToPtr(person.UpdatedAt),
	}
}

package dto

import (
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type GenreFilterRequest struct {
	Search *string `form:"search"`
	Page   int     `form:"page,default=1" binding:"omitempty,gte=1,lte=10000"`
	Limit  int     `form:"limit,default=10" binding:"omitempty,gte=1,lte=100"`
}

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

func (r *CreateGenreRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}

type UpdateGenreRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

func (r *UpdateGenreRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}

type GenreResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func ToGenreResponse(genre *entity.Genre) *GenreResponse {
	if genre == nil {
		return nil
	}

	return &GenreResponse{
		ID:        genre.ID,
		Name:      genre.Name,
		CreatedAt: helper.TimeToPtr(genre.CreatedAt),
		UpdatedAt: helper.TimeToPtr(genre.UpdatedAt),
	}
}

func ToGenreResponseSlice(genres []entity.Genre) []*GenreResponse {
	if len(genres) == 0 {
		return []*GenreResponse{}
	}

	responses := make([]*GenreResponse, 0, len(genres))
	for ind := range genres {
		if resp := ToGenreResponse(&genres[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

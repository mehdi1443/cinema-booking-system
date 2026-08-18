package dto

import (
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type SubmitRatingRequest struct {
	MovieID  *int64 `json:"movie_id,omitempty"`
	CinemaID *int64 `json:"cinema_id,omitempty"`
	Score    int    `json:"score" binding:"required,min=1,max=5"`
}

type RatingUserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RatingResponse struct {
	ID        int64               `json:"id"`
	MovieID   *int64              `json:"movie_id,omitempty"`
	CinemaID  *int64              `json:"cinema_id,omitempty"`
	Score     int                 `json:"score"`
	User      *RatingUserResponse `json:"user,omitempty"`
	CreatedAt *time.Time          `json:"created_at,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty"`
}

func ToRatingResponse(rating *entity.Rating) *RatingResponse {
	if rating == nil {
		return nil
	}

	var userResp *RatingUserResponse
	if rating.User != nil {
		fullName := strings.TrimSpace(rating.User.FirstName + " " + rating.User.LastName)
		if fullName == "" {
			fullName = "کاربر"
		}
		userResp = &RatingUserResponse{
			ID:   rating.User.ID,
			Name: fullName,
		}
	}

	return &RatingResponse{
		ID:        rating.ID,
		MovieID:   rating.MovieID,
		CinemaID:  rating.CinemaID,
		Score:     rating.Score,
		User:      userResp,
		CreatedAt: helper.TimeToPtr(rating.CreatedAt),
		UpdatedAt: helper.TimeToPtr(rating.UpdatedAt),
	}
}

func ToRatingResponseSlice(ratings []entity.Rating) []*RatingResponse {
	if len(ratings) == 0 {
		return []*RatingResponse{}
	}

	responses := make([]*RatingResponse, 0, len(ratings))
	for ind := range ratings {
		if resp := ToRatingResponse(&ratings[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type MoviePersonRoleRequest struct {
	PersonID int64                `json:"person_id" binding:"required"`
	Role     enum.MoviePersonRole `json:"role" binding:"required"`
}

type CreateMovieRequest struct {
	Title           string                   `json:"title" binding:"required,min=1,max=255"`
	Description     string                   `json:"description"`
	DurationMinutes int                      `json:"duration_minutes" binding:"required,gt=0"`
	ReleaseDate     time.Time                `json:"release_date" binding:"required"`
	Country         string                   `json:"country"`
	Status          enum.MovieStatus         `json:"status" binding:"required"`
	GenreIDs        []int64                  `json:"genre_ids"`
	People          []MoviePersonRoleRequest `json:"people"`
}

type UpdateMovieRequest struct {
	Title           string                   `json:"title" binding:"required,min=1,max=255"`
	Description     string                   `json:"description"`
	DurationMinutes int                      `json:"duration_minutes" binding:"required,gt=0"`
	ReleaseDate     time.Time                `json:"release_date" binding:"required"`
	Country         string                   `json:"country"`
	Status          enum.MovieStatus         `json:"status" binding:"required"`
	GenreIDs        []int64                  `json:"genre_ids"`
	People          []MoviePersonRoleRequest `json:"people"`
}

type MovieListResponse struct {
	ID              int64            `json:"id"`
	Title           string           `json:"title"`
	DurationMinutes int              `json:"duration_minutes"`
	ReleaseDate     time.Time        `json:"release_date"`
	PosterURLLow    string           `json:"poster_url_low"`
	Status          enum.MovieStatus `json:"status"`
}

type MovieDetailResponse struct {
	ID              int64                  `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	DurationMinutes int                    `json:"duration_minutes"`
	ReleaseDate     time.Time              `json:"release_date"`
	Country         string                 `json:"country"`
	PosterURLLow    string                 `json:"poster_url_low"`
	PosterURLHigh   string                 `json:"poster_url_high"`
	TrailerURL      string                 `json:"trailer_url"`
	Status          enum.MovieStatus       `json:"status"`
	Genres          []GenreResponse        `json:"genres,omitempty"`
	People          []PersonDetailResponse `json:"people,omitempty"`
	CreatedAt       *time.Time             `json:"created_at,omitempty"`
	UpdatedAt       *time.Time             `json:"updated_at,omitempty"`
}

func ToMovieListResponse(movie *entity.Movie) *MovieListResponse {
	if movie == nil {
		return nil
	}

	return &MovieListResponse{
		ID:              movie.ID,
		Title:           movie.Title,
		DurationMinutes: movie.DurationMinutes,
		ReleaseDate:     movie.ReleaseDate,
		PosterURLLow:    movie.PosterURLLow,
		Status:          movie.Status,
	}
}

func ToMovieListResponseSlice(movies []entity.Movie) []*MovieListResponse {
	if len(movies) == 0 {
		return []*MovieListResponse{}
	}

	responses := make([]*MovieListResponse, 0, len(movies))
	for ind := range movies {
		if resp := ToMovieListResponse(&movies[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

func ToMovieDetailResponse(movie *entity.Movie) *MovieDetailResponse {
	if movie == nil {
		return nil
	}

	var genres []GenreResponse
	if len(movie.Genres) > 0 {
		genres = make([]GenreResponse, 0, len(movie.Genres))
		for _, g := range movie.Genres {
			if resp := ToGenreResponse(&g); resp != nil {
				genres = append(genres, *resp)
			}
		}
	}

	var people []PersonDetailResponse
	if len(movie.People) > 0 {
		people = make([]PersonDetailResponse, 0, len(movie.People))
		for _, p := range movie.People {
			if resp := ToPersonDetailResponse(&p); resp != nil {
				people = append(people, *resp)
			}
		}
	}

	return &MovieDetailResponse{
		ID:              movie.ID,
		Title:           movie.Title,
		Description:     movie.Description,
		DurationMinutes: movie.DurationMinutes,
		ReleaseDate:     movie.ReleaseDate,
		Country:         movie.Country,
		PosterURLLow:    movie.PosterURLLow,
		PosterURLHigh:   movie.PosterURLHigh,
		TrailerURL:      movie.TrailerURL,
		Status:          movie.Status,
		Genres:          genres,
		People:          people,
		CreatedAt:       helper.TimeToPtr(movie.CreatedAt),
		UpdatedAt:       helper.TimeToPtr(movie.UpdatedAt),
	}
}

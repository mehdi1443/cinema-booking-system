package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type MovieFilterRequest struct {
	Status      *enum.MovieStatus `form:"status"`
	GenreID     *int64            `form:"genre_id"`
	Country     *string           `form:"country"`
	Search      *string           `form:"search"`
	MinRating   *float64          `form:"min_rating"`
	ReleaseYear *int              `form:"release_year"`
	SortBy      string            `form:"sort_by"`
	SortDir     string            `form:"sort_dir"`
	Page        int               `form:"page,default=1"`
	Limit       int               `form:"limit,default=10"`
}

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
	Title           *string                   `json:"title" binding:"omitempty,min=1,max=255"`
	Description     *string                   `json:"description"`
	DurationMinutes *int                      `json:"duration_minutes" binding:"omitempty,gt=0"`
	ReleaseDate     *time.Time                `json:"release_date"`
	Country         *string                   `json:"country"`
	Status          *enum.MovieStatus         `json:"status"`
	GenreIDs        *[]int64                  `json:"genre_ids"`
	People          *[]MoviePersonRoleRequest `json:"people"`
}

type MovieGenreResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MoviePersonResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	ImageURLLow string               `json:"image_url_low"`
	Role        enum.MoviePersonRole `json:"role"`
}

type MovieListResponse struct {
	ID              int64            `json:"id"`
	Title           string           `json:"title"`
	DurationMinutes int              `json:"duration_minutes"`
	ReleaseDate     time.Time        `json:"release_date"`
	PosterURLLow    string           `json:"poster_url_low"`
	Status          enum.MovieStatus `json:"status"`
	RatingAvg       float64          `json:"rating_avg"`
	RatingCount     int              `json:"rating_count"`
}

type MovieDetailResponse struct {
	ID              int64                 `json:"id"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	DurationMinutes int                   `json:"duration_minutes"`
	ReleaseDate     time.Time             `json:"release_date"`
	Country         string                `json:"country"`
	PosterURLLow    string                `json:"poster_url_low"`
	PosterURLHigh   string                `json:"poster_url_high"`
	TrailerURL      string                `json:"trailer_url"`
	Status          enum.MovieStatus      `json:"status"`
	RatingAvg       float64               `json:"rating_avg"`
	RatingCount     int                   `json:"rating_count"`
	Genres          []MovieGenreResponse  `json:"genres,omitempty"`
	People          []MoviePersonResponse `json:"people,omitempty"`
	CreatedAt       *time.Time            `json:"created_at,omitempty"`
	UpdatedAt       *time.Time            `json:"updated_at,omitempty"`
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
		RatingAvg:       movie.RatingAvg,
		RatingCount:     movie.RatingCount,
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

	var genres []MovieGenreResponse
	if len(movie.Genres) > 0 {
		genres = make([]MovieGenreResponse, 0, len(movie.Genres))
		for _, g := range movie.Genres {
			genres = append(genres, MovieGenreResponse{
				ID:   g.ID,
				Name: g.Name,
			})
		}
	}

	var people []MoviePersonResponse
	if len(movie.People) > 0 {
		people = make([]MoviePersonResponse, 0, len(movie.People))
		for _, mp := range movie.People {
			people = append(people, MoviePersonResponse{
				ID:          mp.Person.ID,
				Name:        mp.Person.Name,
				ImageURLLow: mp.Person.ImageURLLow,
				Role:        mp.Role,
			})
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
		RatingAvg:       movie.RatingAvg,
		RatingCount:     movie.RatingCount,
		Genres:          genres,
		People:          people,
		CreatedAt:       helper.TimeToPtr(movie.CreatedAt),
		UpdatedAt:       helper.TimeToPtr(movie.UpdatedAt),
	}
}

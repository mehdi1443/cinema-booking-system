package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type Movie struct {
	ID              int64            `json:"id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	DurationMinutes int              `json:"duration_minutes"`
	ReleaseDate     time.Time        `json:"release_date"`
	Country         string           `json:"country"`
	PosterURLLow    string           `json:"poster_url_low"`
	PosterURLHigh   string           `json:"poster_url_high"`
	TrailerURL      string           `json:"trailer_url"`
	Status          enum.MovieStatus `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at"`
	Genres          []Genre          `json:"genres,omitempty"`
	People          []Person         `json:"people,omitempty"`
}

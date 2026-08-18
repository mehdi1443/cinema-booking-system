package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type Movie struct {
	ID              int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string           `gorm:"type:varchar(255);not null" json:"title"`
	Description     string           `gorm:"type:text" json:"description"`
	DurationMinutes int              `gorm:"not null" json:"duration_minutes"`
	ReleaseDate     time.Time        `gorm:"type:date;not null" json:"release_date"`
	Country         string           `gorm:"type:varchar(100)" json:"country"`
	PosterURLLow    string           `gorm:"type:varchar(500)" json:"poster_url_low"`
	PosterURLHigh   string           `gorm:"type:varchar(500)" json:"poster_url_high"`
	TrailerURL      string           `gorm:"type:varchar(500)" json:"trailer_url"`
	Status          enum.MovieStatus `gorm:"type:varchar(50);not null" json:"status"`
	RatingAvg       float64          `gorm:"type:decimal(3,2);not null;default:0.00" json:"rating_avg"`
	RatingCount     int              `gorm:"not null;default:0" json:"rating_count"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time       `gorm:"index" json:"deleted_at,omitempty"`

	Genres []Genre       `gorm:"many2many:movie_genres;" json:"genres,omitempty"`
	People []MoviePerson `gorm:"foreignKey:MovieID" json:"people,omitempty"`
}

package entity

import (
	"time"
)

type Rating struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	MovieID   *int64    `json:"movie_id,omitempty"`
	CinemaID  *int64    `json:"cinema_id,omitempty"`
	Score     int       `gorm:"not null" json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie  *Movie  `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Cinema *Cinema `gorm:"foreignKey:CinemaID" json:"cinema,omitempty"`
}

func (Rating) TableName() string {
	return "ratings"
}

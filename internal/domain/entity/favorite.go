package entity

import "time"

type Favorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null" json:"user_id"`
	MovieID   *int64    `gorm:"column:movie_id" json:"movie_id,omitempty"`
	CinemaID  *int64    `gorm:"column:cinema_id" json:"cinema_id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie  *Movie  `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Cinema *Cinema `gorm:"foreignKey:CinemaID" json:"cinema,omitempty"`
}

package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type Comment struct {
	ID            int64              `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64              `gorm:"not null" json:"user_id"`
	MovieID       *int64             `json:"movie_id,omitempty"`
	CinemaID      *int64             `json:"cinema_id,omitempty"`
	ParentID      *int64             `json:"parent_id,omitempty"`
	Body          string             `gorm:"type:text;not null" json:"body"`
	LikesCount    int                `gorm:"default:0;not null" json:"likes_count"`
	DislikesCount int                `gorm:"default:0;not null" json:"dislikes_count"`
	Status        enum.CommentStatus `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`

	User    *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie   *Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Cinema  *Cinema   `gorm:"foreignKey:CinemaID" json:"cinema,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

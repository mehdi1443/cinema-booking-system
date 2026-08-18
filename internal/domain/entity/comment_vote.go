package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CommentVote struct {
	UserID    int64         `gorm:"primaryKey" json:"user_id"`
	CommentID int64         `gorm:"primaryKey;column:review_id" json:"comment_id"`
	Vote      enum.VoteType `gorm:"type:varchar(50);not null" json:"vote"`
	CreatedAt time.Time     `json:"created_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

func (CommentVote) TableName() string {
	return "review_votes"
}

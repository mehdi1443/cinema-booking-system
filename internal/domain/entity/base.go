package entity

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

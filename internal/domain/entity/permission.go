package entity

import (
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"gorm.io/gorm"
)

type Permission struct {
	BaseEntity
	DeletedAt   gorm.DeletedAt
	Name        string                  `json:"name" gorm:"uniqueIndex;not null"`
	Title       string                  `json:"title"`
	Category    enum.PermissionCategory `json:"category"`
	Description string                  `json:"description"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

package entity

import "gorm.io/gorm"

type Role struct {
	BaseEntity
	DeletedAt   gorm.DeletedAt
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Title       string `json:"title"`
	Description string `json:"description" gorm:"-"`

	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Users       []*User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

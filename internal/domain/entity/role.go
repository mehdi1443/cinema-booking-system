package entity

import "gorm.io/gorm"

type Role struct {
	BaseEntity
	DeletedAt   gorm.DeletedAt
	Name        string        `json:"name"`
	Title       string        `json:"description"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Users       []*User       `json:"users,omitempty" gorm:"many2many:role_users;"`
}

package entity

import "github.com/mehdi1443/cinema-booking-system/internal/domain/enum"

type Permission struct {
	BaseEntity
	Name        string                  `json:"name"`
	Title       string                  `json:"title"`
	Category    enum.PermissionCategory `json:"category"`
	Description string                  `json:"description"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:permission_roles;"`
}

// ? enum : PermissionType
// example :
/*
var p PermissionType
ID 			 Create With DB
Name:        p.String()
Title:       p.Description()
Category:    permType.Category()
Description: permType.Description()
*/

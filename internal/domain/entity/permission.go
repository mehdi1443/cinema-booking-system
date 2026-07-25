package entity

import "github.com/mehdi1443/cinema-booking-system/internal/domain/enum"

type Permission struct {
	BaseEntity
	Name        string                  `json:"name"`
	Title       string                  `json:"title"`
	Category    enum.PermissionCategory `json:"category"`
	Description string                  `json:"description"`
}

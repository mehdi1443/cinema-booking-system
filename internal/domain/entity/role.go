package entity

type Role struct {
	BaseEntity
	Name        string       `json:"name"`
	Title       string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

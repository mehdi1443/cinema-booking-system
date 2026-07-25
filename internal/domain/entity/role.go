package entity

type Role struct {
	BaseEntity
	Name        string       `json:"name"`
	title       string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
}

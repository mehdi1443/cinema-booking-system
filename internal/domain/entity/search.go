package entity

import "github.com/mehdi1443/cinema-booking-system/internal/domain/enum"

type SearchDocument struct {
	ID       string       `json:"id"`
	EntityID int64        `json:"entity_id"`
	Type     enum.DocType `json:"type"`
	Title    string       `json:"title"`
	Subtitle string       `json:"subtitle"`
	ImageURL string       `json:"image_url"`
}

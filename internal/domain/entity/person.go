package entity

import "time"

type Person struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Biography    string    `json:"biography"`
	ImageURLLow  string    `json:"image_url_low"`
	ImageURLHigh string    `json:"image_url_high"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

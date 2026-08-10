package entity

type Cinema struct {
	ID           int64        `json:"id"`
	CityID       int64        `json:"city_id"`
	Name         string       `json:"name"`
	Description  *string      `json:"description,omitempty"`
	Address      string       `json:"address"`
	Latitude     *float64     `json:"latitude,omitempty"`
	Longitude    *float64     `json:"longitude,omitempty"`
	Phone        *string      `json:"phone,omitempty"`
	ImageURLLow  *string      `json:"image_url_low,omitempty"`
	ImageURLHigh *string      `json:"image_url_high,omitempty"`
	Status       CinemaStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`

	City       *City      `json:"city,omitempty"`
	Facilities []Facility `json:"facilities,omitempty"`
}

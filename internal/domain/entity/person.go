package entity

type Person struct {
	BaseEntity
	Name         string `json:"name"`
	Biography    string `json:"biography"`
	ImageURLLow  string `json:"image_url_low"`
	ImageURLHigh string `json:"image_url_high"`
}

package entity

type City struct {
	ID         int64     `json:"id"`
	ProvinceID int64     `json:"province_id"`
	Name       string    `json:"name"`
	Province   *Province `json:"province,omitempty"`
}

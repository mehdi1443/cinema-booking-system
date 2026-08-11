package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type Cinema struct {
	ID           int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	CityID       int64             `gorm:"not null" json:"city_id"`
	Name         string            `gorm:"type:varchar(255);not null" json:"name"`
	Description  *string           `gorm:"type:text" json:"description,omitempty"`
	Address      string            `gorm:"type:text;not null" json:"address"`
	Latitude     *float64          `json:"latitude,omitempty"`
	Longitude    *float64          `json:"longitude,omitempty"`
	Phone        *string           `gorm:"type:varchar(20)" json:"phone,omitempty"`
	ImageURLLow  *string           `gorm:"type:varchar(255)" json:"image_url_low,omitempty"`
	ImageURLHigh *string           `gorm:"type:varchar(255)" json:"image_url_high,omitempty"`
	Status       enum.CinemaStatus `gorm:"type:varchar(50)" json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`

	City       *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Facilities []Facility `gorm:"many2many:cinema_facilities;foreignKey:ID;joinForeignKey:CinemaID;references:ID;joinReferences:FacilityID" json:"facilities,omitempty"`
}

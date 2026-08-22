package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type Cinema struct {
	BaseEntity
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
	RatingAvg    float64           `gorm:"type:decimal(3,2);not null;default:0.00" json:"rating_avg"`
	RatingCount  int               `gorm:"not null;default:0" json:"rating_count"`
	DeletedAt    *time.Time        `gorm:"index" json:"deleted_at,omitempty"`

	City       *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Facilities []Facility `gorm:"many2many:cinema_facilities;foreignKey:ID;joinForeignKey:CinemaID;references:ID;joinReferences:FacilityID" json:"facilities,omitempty"`
}

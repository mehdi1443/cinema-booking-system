package entity

import "github.com/mehdi1443/cinema-booking-system/internal/domain/enum"

type MoviePerson struct {
	MovieID  int64                `gorm:"primaryKey;not null" json:"movie_id"`
	PersonID int64                `gorm:"primaryKey;not null" json:"person_id"`
	Role     enum.MoviePersonRole `gorm:"primaryKey;type:movie_person_role;not null" json:"role"`
	Person   Person               `gorm:"foreignKey:PersonID" json:"person"`
}

func (MoviePerson) TableName() string {
	return "movie_people"
}

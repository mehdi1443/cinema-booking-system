package entity

type Facility struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}

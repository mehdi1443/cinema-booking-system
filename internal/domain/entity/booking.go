package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"gorm.io/gorm"
)

type Hall struct {
	BaseEntity
	CinemaID  int64           `gorm:"not null;index" json:"cinema_id"`
	Name      string          `gorm:"type:varchar(100);not null" json:"name"`
	Status    enum.HallStatus `gorm:"type:hall_status;not null" json:"status"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Cinema *Cinema `gorm:"foreignKey:CinemaID" json:"cinema,omitempty"`
	Seats  []Seat  `gorm:"foreignKey:HallID" json:"seats,omitempty"`
}

type Seat struct {
	BaseEntity
	HallID    int64           `gorm:"not null;index" json:"hall_id"`
	RowLabel  string          `gorm:"type:varchar(10);not null" json:"row_label"`
	Number    int             `gorm:"not null" json:"number"`
	Type      enum.SeatType   `gorm:"column:seat_type;type:seat_type;not null" json:"type"`
	Status    enum.SeatStatus `gorm:"type:seat_status;not null" json:"status"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	Hall *Hall `gorm:"foreignKey:HallID" json:"hall,omitempty"`
}

type Showtime struct {
	BaseEntity
	HallID    int64               `gorm:"not null;index" json:"hall_id"`
	MovieID   int64               `gorm:"not null;index" json:"movie_id"`
	StartsAt  time.Time           `gorm:"not null;index" json:"starts_at"`
	EndsAt    time.Time           `gorm:"not null" json:"ends_at"`
	Status    enum.ShowtimeStatus `gorm:"type:showtime_status;not null" json:"status"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`

	Hall  *Hall          `gorm:"foreignKey:HallID" json:"hall,omitempty"`
	Movie *Movie         `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Seats []ShowtimeSeat `gorm:"foreignKey:ShowtimeID" json:"seats,omitempty"`
}

type Reservation struct {
	BaseEntity
	UserID      int64                  `gorm:"not null;index" json:"user_id"`
	ShowtimeID  int64                  `gorm:"not null;index" json:"showtime_id"`
	Status      enum.ReservationStatus `gorm:"type:reservation_status;not null;index" json:"status"`
	ExpiresAt   time.Time              `gorm:"not null;index" json:"expires_at"`
	TotalAmount int64                  `gorm:"not null" json:"total_amount"`

	User     *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Showtime *Showtime         `gorm:"foreignKey:ShowtimeID" json:"showtime,omitempty"`
	Items    []ReservationItem `gorm:"foreignKey:ReservationID" json:"items,omitempty"`
	Payments []Payment         `gorm:"foreignKey:ReservationID" json:"payments,omitempty"`
}

type ShowtimeSeat struct {
	BaseEntity
	ShowtimeID    int64                   `gorm:"not null;index" json:"showtime_id"`
	SeatID        int64                   `gorm:"not null;index" json:"seat_id"`
	RowLabel      string                  `gorm:"type:varchar(10);not null" json:"row_label"`
	SeatNumber    int                     `gorm:"not null" json:"seat_number"`
	SeatType      enum.SeatType           `gorm:"type:seat_type;not null" json:"seat_type"`
	Price         int64                   `gorm:"not null" json:"price"`
	Status        enum.ShowtimeSeatStatus `gorm:"type:showtime_seat_status;not null;index" json:"status"`
	ReservationID *int64                  `gorm:"index" json:"reservation_id,omitempty"`
	HeldUntil     *time.Time              `gorm:"index" json:"held_until,omitempty"`

	Showtime    *Showtime    `gorm:"foreignKey:ShowtimeID" json:"showtime,omitempty"`
	Seat        *Seat        `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
	Reservation *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
}

type ReservationItem struct {
	BaseEntity
	ReservationID  int64 `gorm:"not null;index" json:"reservation_id"`
	ShowtimeSeatID int64 `gorm:"not null;index" json:"showtime_seat_id"`
	UnitPrice      int64 `gorm:"not null" json:"unit_price"`

	ShowtimeSeat *ShowtimeSeat `gorm:"foreignKey:ShowtimeSeatID" json:"showtime_seat,omitempty"`
}

type Payment struct {
	BaseEntity
	ReservationID  int64                `gorm:"not null;index" json:"reservation_id"`
	UserID         int64                `gorm:"not null;index" json:"user_id"`
	Amount         int64                `gorm:"not null" json:"amount"`
	Provider       enum.PaymentProvider `gorm:"type:payment_provider;not null" json:"provider"`
	Status         enum.PaymentStatus   `gorm:"type:payment_status;not null;index" json:"status"`
	Authority      string               `gorm:"type:varchar(100);not null;uniqueIndex" json:"authority"`
	ReferenceID    *string              `gorm:"type:varchar(100);uniqueIndex" json:"reference_id,omitempty"`
	IdempotencyKey string               `gorm:"type:varchar(100);not null" json:"-"`
	FailureReason  *string              `gorm:"type:varchar(255)" json:"failure_reason,omitempty"`
	PaidAt         *time.Time           `json:"paid_at,omitempty"`

	Reservation *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

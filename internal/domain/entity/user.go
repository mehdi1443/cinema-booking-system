package entity

import (
	"math"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type User struct {
	BaseEntity
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	PasswordHash  string          `json:"password_hash"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	WalletBalance int64           `json:"wallet_balance"`
	Status        enum.UserStatus `json:"status"` // مشکل باهاش دارم چرا اینجاست و چرا باشه و چرا نباشه

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

var (
	MaxWalletBalance int64 = math.MaxInt64
)

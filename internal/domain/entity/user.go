package entity

import (
	"math"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"gorm.io/gorm"
)

type User struct {
	BaseEntity
	DeletedAt     gorm.DeletedAt
	Phone         string          `json:"phone"`
	Email         string          `json:"email"`
	PasswordHash  string          `json:"password_hash"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	WalletBalance int64           `json:"wallet_balance"`
	Status        enum.UserStatus `json:"status"` // مشکل باهاش دارم چرا اینجاست و چرا باشه و چرا نباشه

	IsPhoneVerified bool `json:"is_verified_phone"`
	IsEmailVerified bool `json:"is_verified_email"`

	PhoneVerifiedAt time.Time `json:"phone_verified_time"`
	EmailVerifiedAt time.Time `json:"email_verified_time"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

var (
	MaxWalletBalance int64 = math.MaxInt64
)

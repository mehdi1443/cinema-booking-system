package entity

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"gorm.io/gorm"
)

type User struct {
	BaseEntity
	DeletedAt     gorm.DeletedAt
	Phone         string          `json:"phone" gorm:"uniqueIndex;not null"`
	Email         *string         `json:"email" gorm:"uniqueIndex"`
	PasswordHash  *string         `json:"password_hash"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	WalletBalance int64           `json:"wallet_balance"`
	Status        enum.UserStatus `json:"status"`

	IsPhoneVerified bool `json:"is_verified_phone"`
	IsEmailVerified bool `json:"is_verified_email"`

	PhoneVerifiedAt *time.Time `json:"phone_verified_time"`
	EmailVerifiedAt *time.Time `json:"email_verified_time"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

const MaxWalletBalance int64 = 10_000_000_000

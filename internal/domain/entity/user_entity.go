package entity

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusBlocked  UserStatus = "BLOCKED"
)

type User struct {
	BaseEntity
	Phone         string
	Email         *string
	PasswordHash  *string
	FirstName     *string
	LastName      *string
	WalletBalance int64
	Status        UserStatus
}

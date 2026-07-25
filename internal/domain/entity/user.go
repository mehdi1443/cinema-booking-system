package entity

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusBlocked  UserStatus = "BLOCKED"
)

type User struct {
	BaseEntity
	Phone         string     `json:"phone"`
	Email         *string    `json:"email"`
	PasswordHash  *string    `json:"-"`
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	WalletBalance int64      `json:"wallet_balance"`
	Status        UserStatus `json:"status"`

	Roles []Role `json:"roles,omitempty"`
}

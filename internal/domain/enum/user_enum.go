package enum

import "strings"

type UserStatus string

const (
	UserActive   UserStatus = "ACTIVE"
	UserInactive UserStatus = "INACTIVE"
	UserBlocked  UserStatus = "BLOCKED"
)

func (s UserStatus) String() string {
	return string(s)
}

func AllUserStatuses() []UserStatus {
	return []UserStatus{
		UserActive,
		UserInactive,
		UserBlocked,
	}
}

func (s UserStatus) IsValid() bool {
	for _, status := range AllUserStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

func ParseUserStatus(str string) (UserStatus, bool) {
	parsed := UserStatus(strings.ToUpper(strings.TrimSpace(str)))

	if parsed.IsValid() {
		return parsed, true
	}

	return "", false
}

package enum

import "strings"

type CinemaStatus string

const (
	CinemaStatusActive   CinemaStatus = "ACTIVE"
	CinemaStatusInactive CinemaStatus = "INACTIVE"
)

func (s CinemaStatus) String() string {
	return string(s)
}

func AllCinemaStatuses() []CinemaStatus {
	return []CinemaStatus{
		CinemaStatusActive,
		CinemaStatusInactive,
	}
}

func (s CinemaStatus) IsValid() bool {
	for _, status := range AllCinemaStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

func ParseCinemaStatus(str string) (CinemaStatus, bool) {
	parsed := CinemaStatus(strings.ToUpper(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}

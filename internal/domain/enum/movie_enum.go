package enum

import "strings"

type MovieStatus string

const (
	MovieUpcoming   MovieStatus = "UPCOMING"
	MovieNowShowing MovieStatus = "NOW_SHOWING"
	MovieEnded      MovieStatus = "ENDED"
	MovieInactive   MovieStatus = "INACTIVE"
)

func (s MovieStatus) String() string {
	return string(s)
}

func AllMovieStatuses() []MovieStatus {
	return []MovieStatus{
		MovieUpcoming,
		MovieNowShowing,
		MovieEnded,
		MovieInactive,
	}
}

func (s MovieStatus) IsValid() bool {
	for _, status := range AllMovieStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

func ParseMovieStatus(str string) (MovieStatus, bool) {
	parsed := MovieStatus(strings.ToUpper(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}

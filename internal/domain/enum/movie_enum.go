package enum

type MovieStatus string

const (
	MovieUpcoming   MovieStatus = "UPCOMING"
	MovieNowShowing MovieStatus = "NOW_SHOWING"
	MovieEnded      MovieStatus = "ENDED"
	MovieInactive   MovieStatus = "INACTIVE"
)

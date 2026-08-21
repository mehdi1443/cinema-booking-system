package sort

import "strings"

type MovieSortBy string

const (
	MovieSortByCreatedAt       MovieSortBy = "created_at"
	MovieSortByTitle           MovieSortBy = "title"
	MovieSortByReleaseDate     MovieSortBy = "release_date"
	MovieSortByRatingAvg       MovieSortBy = "rating_avg"
	MovieSortByDurationMinutes MovieSortBy = "duration_minutes"
)

func (s MovieSortBy) String() string {
	return string(s)
}

func (s MovieSortBy) DBColumn() string {
	return string(s)
}

func AllMovieSortBys() []MovieSortBy {
	return []MovieSortBy{
		MovieSortByCreatedAt,
		MovieSortByTitle,
		MovieSortByReleaseDate,
		MovieSortByRatingAvg,
		MovieSortByDurationMinutes,
	}
}

func (s MovieSortBy) IsValid() bool {
	for _, valid := range AllMovieSortBys() {
		if s == valid {
			return true
		}
	}
	return false
}

func ParseMovieSortBy(str string) (MovieSortBy, bool) {
	parsed := MovieSortBy(strings.ToLower(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}

type SortDirection string

const (
	SortAscending  SortDirection = "ASC"
	SortDescending SortDirection = "DESC"
)

func (d SortDirection) String() string {
	return string(d)
}

func AllSortDirections() []SortDirection {
	return []SortDirection{
		SortAscending,
		SortDescending,
	}
}

func (d SortDirection) IsValid() bool {
	for _, valid := range AllSortDirections() {
		if d == valid {
			return true
		}
	}
	return false
}

func ParseSortDirection(str string) (SortDirection, bool) {
	parsed := SortDirection(strings.ToUpper(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}

package sort

import "strings"

type CinemaSortBy string

const (
	CinemaSortByCreatedAt CinemaSortBy = "created_at"
	CinemaSortByName      CinemaSortBy = "name"
	CinemaSortByRatingAvg CinemaSortBy = "rating_avg"
)

func (s CinemaSortBy) String() string {
	return string(s)
}

func (s CinemaSortBy) DBColumn() string {
	return string(s)
}

func AllCinemaSortBys() []CinemaSortBy {
	return []CinemaSortBy{
		CinemaSortByCreatedAt,
		CinemaSortByName,
		CinemaSortByRatingAvg,
	}
}

func (s CinemaSortBy) IsValid() bool {
	for _, valid := range AllCinemaSortBys() {
		if s == valid {
			return true
		}
	}
	return false
}

func ParseCinemaSortBy(str string) (CinemaSortBy, bool) {
	parsed := CinemaSortBy(strings.ToLower(strings.TrimSpace(str)))
	if parsed.IsValid() {
		return parsed, true
	}
	return "", false
}

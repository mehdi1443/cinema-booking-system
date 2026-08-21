package enum

import "strings"

type BucketType string

const (
	MoviePoster  BucketType = "moviePoster"
	MovieTrailer BucketType = "movieTrailer"
	ActorPhoto   BucketType = "actorPhoto"
	CinemaPhoto  BucketType = "cinemaPhoto"
)

func (bt BucketType) String() string {
	return string(bt)
}

func AllBucketTypes() []BucketType {
	return []BucketType{
		MoviePoster,
		MovieTrailer,
		ActorPhoto,
		CinemaPhoto,
	}
}

func (bt BucketType) IsValid() bool {
	for _, b := range AllBucketTypes() {
		if bt == b {
			return true
		}
	}
	return false
}

func ParseBucketType(str string) (BucketType, bool) {
	str = strings.TrimSpace(str)
	for _, b := range AllBucketTypes() {
		if strings.EqualFold(string(b), str) {
			return b, true
		}
	}
	return "", false
}

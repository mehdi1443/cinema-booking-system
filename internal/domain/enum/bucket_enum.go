package enum

type BucketType uint

const (
	MoviePoster BucketType = iota + 1
	MovieTrailer
	ActorPhoto
	CinemaPhoto
)

func (bt BucketType) String() string {
	switch bt {
	case MoviePoster:
		return "moviePoster"
	case MovieTrailer:
		return "movieTrailer"
	case ActorPhoto:
		return "actorPhoto"
	case CinemaPhoto:
		return "cinemaPhoto"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		MoviePoster,
		MovieTrailer,
		ActorPhoto,
		CinemaPhoto,
	}
}

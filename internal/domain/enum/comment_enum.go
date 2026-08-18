package enum

type CommentStatus string

const (
	CommentStatusPending   CommentStatus = "PENDING"
	CommentStatusPublished CommentStatus = "PUBLISHED"
	CommentStatusRejected  CommentStatus = "REJECTED"
	CommentStatusSpoiler   CommentStatus = "SPOILER"
)

func (s CommentStatus) IsValid() bool {
	switch s {
	case CommentStatusPending, CommentStatusPublished, CommentStatusRejected, CommentStatusSpoiler:
		return true
	}
	return false
}

type VoteType string

const (
	VoteTypeLike    VoteType = "LIKE"
	VoteTypeDislike VoteType = "DISLIKE"
)

func (v VoteType) IsValid() bool {
	switch v {
	case VoteTypeLike, VoteTypeDislike:
		return true
	}
	return false
}

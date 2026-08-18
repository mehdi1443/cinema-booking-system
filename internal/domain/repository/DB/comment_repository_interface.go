package DB

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CommentInterface interface {
	Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	GetByID(ctx context.Context, id int64) (*entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	UpdateStatus(ctx context.Context, id int64, status enum.CommentStatus) error
	Delete(ctx context.Context, id int64) error

	GetByMovieID(ctx context.Context, movieID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error)
	GetByCinemaID(ctx context.Context, cinemaID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error)
	GetRepliesByParentID(ctx context.Context, parentID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error)

	AddVote(ctx context.Context, vote *entity.CommentVote) error
	RemoveVote(ctx context.Context, userID int64, commentID int64) error
	GetUserVote(ctx context.Context, userID int64, commentID int64) (*entity.CommentVote, error)

	GetListAdmin(ctx context.Context, status *enum.CommentStatus, limit, offset int) ([]entity.Comment, error)
}

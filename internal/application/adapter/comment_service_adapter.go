package service_contract

import (
	"context"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CommentServiceInterface interface {
	Create(ctx context.Context, userID int64, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.CommentResponse, error)
	GetByMovieID(ctx context.Context, movieID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error)
	GetByMovieIDAdmin(ctx context.Context, movieID int64, page, limit int) ([]*dto.CommentResponse, error)
	GetByCinemaID(ctx context.Context, cinemaID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error)
	GetByCinemaIDAdmin(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.CommentResponse, error)
	GetReplies(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error)
	GetRepliesAdmin(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error)
	Update(ctx context.Context, userID int64, commentID int64, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	UpdateStatus(ctx context.Context, commentID int64, status enum.CommentStatus) error
	Delete(ctx context.Context, userID int64, commentID int64) error

	Vote(ctx context.Context, userID int64, commentID int64, vote enum.VoteType) error
	RemoveVote(ctx context.Context, userID int64, commentID int64) error

	GetListAdmin(ctx context.Context, status *enum.CommentStatus, page, limit int) ([]*dto.CommentResponse, error)
}

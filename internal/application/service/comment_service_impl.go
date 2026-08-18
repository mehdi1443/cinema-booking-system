package service

import (
	"context"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type commentService struct {
	commentRepo repository_contract.CommentInterface
}

func NewCommentService(commentRepo repository_contract.CommentInterface) service_contract.CommentServiceInterface {
	return &commentService{
		commentRepo: commentRepo,
	}
}

func (s *commentService) Create(ctx context.Context, userID int64, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if req.MovieID == nil && req.CinemaID == nil && req.ParentID == nil {
		return nil, exception.ErrInvalidInput
	}

	comment := &entity.Comment{
		UserID:   userID,
		MovieID:  req.MovieID,
		CinemaID: req.CinemaID,
		ParentID: req.ParentID,
		Body:     req.Body,
		Status:   enum.CommentStatusPending,
	}

	createdComment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return dto.ToCommentResponse(createdComment, nil), nil
}

func (s *commentService) GetByID(ctx context.Context, id int64) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToCommentResponse(comment, nil), nil
}

func (s *commentService) GetByMovieID(ctx context.Context, movieID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	// ارسال هم‌زمان وضعیت‌های Published و Spoiler
	comments, err := s.commentRepo.GetByMovieID(ctx, movieID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, err
	}

	return s.enrichCommentsWithVotes(ctx, comments, currentUserID)
}

func (s *commentService) GetByMovieIDAdmin(ctx context.Context, movieID int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByMovieID(ctx, movieID, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

func (s *commentService) GetByCinemaID(ctx context.Context, cinemaID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	// ارسال هم‌زمان وضعیت‌های Published و Spoiler
	comments, err := s.commentRepo.GetByCinemaID(ctx, cinemaID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, err
	}

	return s.enrichCommentsWithVotes(ctx, comments, currentUserID)
}

func (s *commentService) GetByCinemaIDAdmin(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByCinemaID(ctx, cinemaID, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

func (s *commentService) GetReplies(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	// ارسال هم‌زمان وضعیت‌های Published و Spoiler
	replies, err := s.commentRepo.GetRepliesByParentID(ctx, parentID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, err
	}

	return dto.ToCommentResponseSlice(replies), nil
}

func (s *commentService) GetRepliesAdmin(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	replies, err := s.commentRepo.GetRepliesByParentID(ctx, parentID, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.ToCommentResponseSlice(replies), nil
}

func (s *commentService) Update(ctx context.Context, userID int64, commentID int64, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID {
		return nil, exception.ErrUnauthorized
	}

	comment.Body = req.Body

	updatedComment, err := s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return dto.ToCommentResponse(updatedComment, nil), nil
}

func (s *commentService) UpdateStatus(ctx context.Context, commentID int64, status enum.CommentStatus) error {
	return s.commentRepo.UpdateStatus(ctx, commentID, status)
}

func (s *commentService) Delete(ctx context.Context, userID int64, commentID int64) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return exception.ErrUnauthorized
	}

	return s.commentRepo.Delete(ctx, commentID)
}

func (s *commentService) Vote(ctx context.Context, userID int64, commentID int64, vote enum.VoteType) error {
	_, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	commentVote := &entity.CommentVote{
		UserID:    userID,
		CommentID: commentID,
		Vote:      vote,
	}

	return s.commentRepo.AddVote(ctx, commentVote)
}

func (s *commentService) RemoveVote(ctx context.Context, userID int64, commentID int64) error {
	return s.commentRepo.RemoveVote(ctx, userID, commentID)
}

func normalizePagination(page, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func (s *commentService) enrichCommentsWithVotes(ctx context.Context, comments []entity.Comment, currentUserID *int64) ([]*dto.CommentResponse, error) {
	responses := make([]*dto.CommentResponse, 0, len(comments))
	for i := range comments {
		var userVote *enum.VoteType
		if currentUserID != nil {
			vote, err := s.commentRepo.GetUserVote(ctx, *currentUserID, comments[i].ID)
			if err == nil && vote != nil {
				userVote = &vote.Vote
			}
		}
		responses = append(responses, dto.ToCommentResponse(&comments[i], userVote))
	}
	return responses, nil
}

func (s *commentService) GetListAdmin(ctx context.Context, status *enum.CommentStatus, page, limit int) ([]*dto.CommentResponse, error) {
	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetListAdmin(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

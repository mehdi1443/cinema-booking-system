package service

import (
	"context"
	"fmt"
	"strings"

	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	repository_contract "github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
)

type commentService struct {
	commentRepo repository_contract.CommentInterface
	log         domainlogger.Logger
}

func NewCommentService(commentRepo repository_contract.CommentInterface, log domainlogger.Logger) service_contract.CommentServiceInterface {
	return &commentService{
		commentRepo: commentRepo,
		log:         log,
	}
}

func (s *commentService) Create(ctx context.Context, userID int64, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if req == nil || userID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	if req.MovieID == nil && req.CinemaID == nil && req.ParentID == nil {
		return nil, exception.ErrInvalidInput
	}

	trimmedBody := strings.TrimSpace(req.Body)
	if trimmedBody == "" {
		return nil, exception.ErrInvalidInput
	}

	comment := &entity.Comment{
		UserID:   userID,
		MovieID:  req.MovieID,
		CinemaID: req.CinemaID,
		ParentID: req.ParentID,
		Body:     trimmedBody,
		Status:   enum.CommentStatusPending,
	}

	createdComment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	s.log.Info(ctx, "Comment created successfully", domainlogger.Int64("comment_id", createdComment.ID), domainlogger.Int64("user_id", userID))

	return dto.ToCommentResponse(createdComment, nil), nil
}

func (s *commentService) GetByID(ctx context.Context, id int64) (*dto.CommentResponse, error) {
	if id <= 0 {
		return nil, exception.ErrInvalidInput
	}

	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get comment by id %d: %w", id, err)
	}

	return dto.ToCommentResponse(comment, nil), nil
}

func (s *commentService) GetByMovieID(ctx context.Context, movieID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error) {
	if movieID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByMovieID(ctx, movieID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, fmt.Errorf("get comments by movie %d: %w", movieID, err)
	}

	return s.enrichCommentsWithVotes(ctx, comments, currentUserID)
}

func (s *commentService) GetByMovieIDAdmin(ctx context.Context, movieID int64, page, limit int) ([]*dto.CommentResponse, error) {
	if movieID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByMovieID(ctx, movieID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("admin get comments by movie %d: %w", movieID, err)
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

func (s *commentService) GetByCinemaID(ctx context.Context, cinemaID int64, currentUserID *int64, page, limit int) ([]*dto.CommentResponse, error) {
	if cinemaID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByCinemaID(ctx, cinemaID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, fmt.Errorf("get comments by cinema %d: %w", cinemaID, err)
	}

	return s.enrichCommentsWithVotes(ctx, comments, currentUserID)
}

func (s *commentService) GetByCinemaIDAdmin(ctx context.Context, cinemaID int64, page, limit int) ([]*dto.CommentResponse, error) {
	if cinemaID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	comments, err := s.commentRepo.GetByCinemaID(ctx, cinemaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("admin get comments by cinema %d: %w", cinemaID, err)
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

func (s *commentService) GetReplies(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error) {
	if parentID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	replies, err := s.commentRepo.GetRepliesByParentID(ctx, parentID, limit, offset, enum.CommentStatusPublished, enum.CommentStatusSpoiler)
	if err != nil {
		return nil, fmt.Errorf("get comment replies %d: %w", parentID, err)
	}

	return dto.ToCommentResponseSlice(replies), nil
}

func (s *commentService) GetRepliesAdmin(ctx context.Context, parentID int64, page, limit int) ([]*dto.CommentResponse, error) {
	if parentID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	page, limit, offset := normalizePagination(page, limit)

	replies, err := s.commentRepo.GetRepliesByParentID(ctx, parentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("admin get comment replies %d: %w", parentID, err)
	}

	return dto.ToCommentResponseSlice(replies), nil
}

func (s *commentService) Update(ctx context.Context, userID int64, commentID int64, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	if req == nil || userID <= 0 || commentID <= 0 {
		return nil, exception.ErrInvalidInput
	}

	trimmedBody := strings.TrimSpace(req.Body)
	if trimmedBody == "" {
		return nil, exception.ErrInvalidInput
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("get comment before update %d: %w", commentID, err)
	}

	if comment.UserID != userID {
		return nil, exception.ErrNotAccessible
	}

	comment.Body = trimmedBody

	updatedComment, err := s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("update comment %d: %w", commentID, err)
	}

	s.log.Info(ctx, "Comment updated successfully", domainlogger.Int64("comment_id", updatedComment.ID))

	return dto.ToCommentResponse(updatedComment, nil), nil
}

func (s *commentService) UpdateStatus(ctx context.Context, commentID int64, status enum.CommentStatus) error {
	if commentID <= 0 || !status.IsValid() {
		return exception.ErrInvalidInput
	}

	if err := s.commentRepo.UpdateStatus(ctx, commentID, status); err != nil {
		return fmt.Errorf("update comment status %d: %w", commentID, err)
	}

	s.log.Info(ctx, "Comment status updated", domainlogger.Int64("comment_id", commentID), domainlogger.String("status", string(status)))

	return nil
}

func (s *commentService) Delete(ctx context.Context, userID int64, commentID int64) error {
	if userID <= 0 || commentID <= 0 {
		return exception.ErrInvalidInput
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment before delete %d: %w", commentID, err)
	}

	if comment.UserID != userID {
		return exception.ErrNotAccessible
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("delete comment %d: %w", commentID, err)
	}

	s.log.Info(ctx, "Comment deleted successfully", domainlogger.Int64("comment_id", commentID))

	return nil
}

func (s *commentService) Vote(ctx context.Context, userID int64, commentID int64, vote enum.VoteType) error {
	if userID <= 0 || commentID <= 0 {
		return exception.ErrInvalidInput
	}

	_, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment before vote %d: %w", commentID, err)
	}

	commentVote := &entity.CommentVote{
		UserID:    userID,
		CommentID: commentID,
		Vote:      vote,
	}

	if err := s.commentRepo.AddVote(ctx, commentVote); err != nil {
		return fmt.Errorf("add vote to comment %d: %w", commentID, err)
	}

	s.log.Info(ctx, "Comment voted", domainlogger.Int64("comment_id", commentID), domainlogger.Int64("user_id", userID), domainlogger.String("vote_type", string(vote)))

	return nil
}

func (s *commentService) RemoveVote(ctx context.Context, userID int64, commentID int64) error {
	if userID <= 0 || commentID <= 0 {
		return exception.ErrInvalidInput
	}

	if err := s.commentRepo.RemoveVote(ctx, userID, commentID); err != nil {
		return fmt.Errorf("remove vote from comment %d: %w", commentID, err)
	}

	s.log.Info(ctx, "Comment vote removed", domainlogger.Int64("comment_id", commentID), domainlogger.Int64("user_id", userID))

	return nil
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
		if currentUserID != nil && *currentUserID > 0 {
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
		return nil, fmt.Errorf("admin get comment list: %w", err)
	}

	return s.enrichCommentsWithVotes(ctx, comments, nil)
}

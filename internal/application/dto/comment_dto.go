package dto

import (
	"strings"
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/application/helper"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
)

type CreateCommentRequest struct {
	MovieID  *int64 `json:"movie_id,omitempty"`
	CinemaID *int64 `json:"cinema_id,omitempty"`
	ParentID *int64 `json:"parent_id,omitempty"`
	Body     string `json:"body" binding:"required,min=3"`
}

type UpdateCommentRequest struct {
	Body string `json:"body" binding:"required,min=3"`
}

type UpdateCommentStatusRequest struct {
	Status enum.CommentStatus `json:"status" binding:"required"`
}

type VoteCommentRequest struct {
	Vote enum.VoteType `json:"vote" binding:"required"`
}

type CommentUserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CommentResponse struct {
	ID              int64                `json:"id"`
	MovieID         *int64               `json:"movie_id,omitempty"`
	CinemaID        *int64               `json:"cinema_id,omitempty"`
	ParentID        *int64               `json:"parent_id,omitempty"`
	Body            string               `json:"body"`
	LikesCount      int                  `json:"likes_count"`
	DislikesCount   int                  `json:"dislikes_count"`
	Status          enum.CommentStatus   `json:"status"`
	User            *CommentUserResponse `json:"user,omitempty"`
	CurrentUserVote *enum.VoteType       `json:"current_user_vote,omitempty"`
	Replies         []CommentResponse    `json:"replies,omitempty"`
	CreatedAt       *time.Time           `json:"created_at,omitempty"`
	UpdatedAt       *time.Time           `json:"updated_at,omitempty"`
}

func ToCommentResponse(comment *entity.Comment, currentUserVote *enum.VoteType) *CommentResponse {
	if comment == nil {
		return nil
	}

	var userResp *CommentUserResponse
	if comment.User != nil {
		fullName := strings.TrimSpace(comment.User.FirstName + " " + comment.User.LastName)
		if fullName == "" {
			fullName = "کاربر"
		}
		userResp = &CommentUserResponse{
			ID:   comment.User.ID,
			Name: fullName,
		}
	}

	var replies []CommentResponse
	if len(comment.Replies) > 0 {
		replies = make([]CommentResponse, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			if r := ToCommentResponse(&reply, nil); r != nil {
				replies = append(replies, *r)
			}
		}
	}

	return &CommentResponse{
		ID:              comment.ID,
		MovieID:         comment.MovieID,
		CinemaID:        comment.CinemaID,
		ParentID:        comment.ParentID,
		Body:            comment.Body,
		LikesCount:      comment.LikesCount,
		DislikesCount:   comment.DislikesCount,
		Status:          comment.Status,
		User:            userResp,
		CurrentUserVote: currentUserVote,
		Replies:         replies,
		CreatedAt:       helper.TimeToPtr(comment.CreatedAt),
		UpdatedAt:       helper.TimeToPtr(comment.UpdatedAt),
	}
}

func ToCommentResponseSlice(comments []entity.Comment) []*CommentResponse {
	if len(comments) == 0 {
		return []*CommentResponse{}
	}

	responses := make([]*CommentResponse, 0, len(comments))
	for ind := range comments {
		if resp := ToCommentResponse(&comments[ind], nil); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}

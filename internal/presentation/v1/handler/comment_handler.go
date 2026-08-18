package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type CommentHandler struct {
	commentService service_contract.CommentServiceInterface
}

func NewCommentHandler(commentService service_contract.CommentServiceInterface) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// POST /comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	comment, err := h.commentService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] CreateComment failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrInvalidInput):
			response.Error(c, http.StatusBadRequest, "At least one target (movie_id, cinema_id, parent_id) must be specified", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Comment created successfully", comment)
}

// GET /comments/:id
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	comment, err := h.commentService.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCommentByID failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "Comment not found", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Comment retrieved successfully", comment)
}

// GET /movies/:id/comments
func (h *CommentHandler) GetMovieComments(c *gin.Context) {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var currentUserID *int64
	if userIDVal, exists := c.Get("user_id"); exists {
		uid := userIDVal.(int64)
		currentUserID = &uid
	}

	comments, err := h.commentService.GetByMovieID(c.Request.Context(), movieID, currentUserID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetMovieComments failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Movie comments retrieved successfully", comments)
}

// GET /admin/movies/:id/comments
func (h *CommentHandler) GetMovieCommentsAdmin(c *gin.Context) {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, err := h.commentService.GetByMovieIDAdmin(c.Request.Context(), movieID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetMovieCommentsAdmin failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Admin movie comments retrieved successfully", comments)
}

// GET /cinemas/:id/comments
func (h *CommentHandler) GetCinemaComments(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var currentUserID *int64
	if userIDVal, exists := c.Get("user_id"); exists {
		uid := userIDVal.(int64)
		currentUserID = &uid
	}

	comments, err := h.commentService.GetByCinemaID(c.Request.Context(), cinemaID, currentUserID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemaComments failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Cinema comments retrieved successfully", comments)
}

// GET /admin/cinemas/:id/comments
func (h *CommentHandler) GetCinemaCommentsAdmin(c *gin.Context) {
	cinemaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid cinema ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, err := h.commentService.GetByCinemaIDAdmin(c.Request.Context(), cinemaID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCinemaCommentsAdmin failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Admin cinema comments retrieved successfully", comments)
}

// GET /comments/:id/replies
func (h *CommentHandler) GetCommentReplies(c *gin.Context) {
	parentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid parent comment ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	replies, err := h.commentService.GetReplies(c.Request.Context(), parentID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCommentReplies failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Comment replies retrieved successfully", replies)
}

// GET /admin/comments/:id/replies
func (h *CommentHandler) GetCommentRepliesAdmin(c *gin.Context) {
	parentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid parent comment ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	replies, err := h.commentService.GetRepliesAdmin(c.Request.Context(), parentID, page, limit)
	if err != nil {
		log.Printf("[SERVICE ERROR] GetCommentRepliesAdmin failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Admin comment replies retrieved successfully", replies)
}

// PUT /comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedComment, err := h.commentService.Update(c.Request.Context(), userID, commentID, &req)
	if err != nil {
		log.Printf("[SERVICE ERROR] UpdateComment failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "Comment not found", err)
		case errors.Is(err, exception.ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "You are not allowed to update this comment", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Comment updated successfully", updatedComment)
}

// PATCH /comments/:id/status
func (h *CommentHandler) UpdateCommentStatus(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.commentService.UpdateStatus(c.Request.Context(), commentID, req.Status); err != nil {
		log.Printf("[SERVICE ERROR] UpdateCommentStatus failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "Comment not found", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Comment status updated successfully", nil)
}

// DELETE /comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	if err := h.commentService.Delete(c.Request.Context(), userID, commentID); err != nil {
		log.Printf("[SERVICE ERROR] DeleteComment failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "Comment not found", err)
		case errors.Is(err, exception.ErrUnauthorized):
			response.Error(c, http.StatusForbidden, "You are not allowed to delete this comment", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Comment deleted successfully", nil)
}

// POST /comments/:id/vote
func (h *CommentHandler) VoteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	var req dto.VoteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.commentService.Vote(c.Request.Context(), userID, commentID, req.Vote); err != nil {
		log.Printf("[SERVICE ERROR] VoteComment failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "Comment not found", err)
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Vote registered successfully", nil)
}

// DELETE /comments/:id/vote
func (h *CommentHandler) RemoveVote(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userID := userIDVal.(int64)

	if err := h.commentService.RemoveVote(c.Request.Context(), userID, commentID); err != nil {
		log.Printf("[SERVICE ERROR] RemoveVote failed: %v", err)
		switch {
		case errors.Is(err, exception.ErrDatabase):
			response.Error(c, http.StatusInternalServerError, "Internal database error occurred", err)
		default:
			response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "Vote removed successfully", nil)
}

func (h *CommentHandler) GetCommentsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var statusPtr *enum.CommentStatus
	if statusStr := c.Query("status"); statusStr != "" {
		st := enum.CommentStatus(statusStr)
		statusPtr = &st
	}

	comments, err := h.commentService.GetListAdmin(c.Request.Context(), statusPtr, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/pkg/response"
)

type CommentHandler struct {
	commentService service_contract.CommentServiceInterface
	log            domainlogger.Logger
}

func NewCommentHandler(commentService service_contract.CommentServiceInterface, log domainlogger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		log:            log,
	}
}

// -----------------------------------------------------------------------------
// Endpoints
// -----------------------------------------------------------------------------

// POST /comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	comment, err := h.commentService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusCreated, "Comment created successfully", comment)
}

// GET /comments/:id
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	comment, err := h.commentService.GetByID(c.Request.Context(), id)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Comment retrieved successfully", comment)
}

// GET /movies/:id/comments
func (h *CommentHandler) GetMovieComments(c *gin.Context) {
	movieID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)
	currentUserID := getOptionalAuthUserID(c)

	comments, err := h.commentService.GetByMovieID(c.Request.Context(), movieID, currentUserID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Movie comments retrieved successfully", comments)
}

// GET /admin/movies/:id/comments
func (h *CommentHandler) GetMovieCommentsAdmin(c *gin.Context) {
	movieID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	comments, err := h.commentService.GetByMovieIDAdmin(c.Request.Context(), movieID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Admin movie comments retrieved successfully", comments)
}

// GET /cinemas/:id/comments
func (h *CommentHandler) GetCinemaComments(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)
	currentUserID := getOptionalAuthUserID(c)

	comments, err := h.commentService.GetByCinemaID(c.Request.Context(), cinemaID, currentUserID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Cinema comments retrieved successfully", comments)
}

// GET /admin/cinemas/:id/comments
func (h *CommentHandler) GetCinemaCommentsAdmin(c *gin.Context) {
	cinemaID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	comments, err := h.commentService.GetByCinemaIDAdmin(c.Request.Context(), cinemaID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Admin cinema comments retrieved successfully", comments)
}

// GET /comments/:id/replies
func (h *CommentHandler) GetCommentReplies(c *gin.Context) {
	parentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	replies, err := h.commentService.GetReplies(c.Request.Context(), parentID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Comment replies retrieved successfully", replies)
}

// GET /admin/comments/:id/replies
func (h *CommentHandler) GetCommentRepliesAdmin(c *gin.Context) {
	parentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	replies, err := h.commentService.GetRepliesAdmin(c.Request.Context(), parentID, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Admin comment replies retrieved successfully", replies)
}

// PUT /comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	updatedComment, err := h.commentService.Update(c.Request.Context(), userID, commentID, &req)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Comment updated successfully", updatedComment)
}

// PATCH /comments/:id/status
func (h *CommentHandler) UpdateCommentStatus(c *gin.Context) {
	commentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if err := h.commentService.UpdateStatus(c.Request.Context(), commentID, req.Status); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Comment status updated successfully", nil)
}

// DELETE /comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	if err := h.commentService.Delete(c.Request.Context(), userID, commentID); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Comment deleted successfully", nil)
}

// POST /comments/:id/vote
func (h *CommentHandler) VoteComment(c *gin.Context) {
	commentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	var req dto.VoteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "Invalid request payload")
		return
	}

	if err := h.commentService.Vote(c.Request.Context(), userID, commentID, req.Vote); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Vote registered successfully", nil)
}

// DELETE /comments/:id/vote
func (h *CommentHandler) RemoveVote(c *gin.Context) {
	commentID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	userID, ok := getAuthUserID(c, h.log)
	if !ok {
		return
	}

	if err := h.commentService.RemoveVote(c.Request.Context(), userID, commentID); err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Vote removed successfully", nil)
}

// GET /admin/comments
func (h *CommentHandler) GetCommentsAdmin(c *gin.Context) {
	page, limit := parsePagination(c)

	var statusPtr *enum.CommentStatus
	if statusStr := c.Query("status"); statusStr != "" {
		st := enum.CommentStatus(statusStr)
		statusPtr = &st
	}

	comments, err := h.commentService.GetListAdmin(c.Request.Context(), statusPtr, page, limit)
	if err != nil {
		Handle(c, err, h.log)
		return
	}

	response.Success(c, http.StatusOK, "Admin comments retrieved successfully", comments)
}

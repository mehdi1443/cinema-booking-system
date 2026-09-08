package router

import (
	"github.com/gin-gonic/gin"

	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
)

func SetUpCommentRouter(rg *gin.RouterGroup, commentHandler *handler.CommentHandler) {
	rg.GET("/comments/:id", commentHandler.GetCommentByID)
	rg.GET("/comments/:id/replies", commentHandler.GetCommentReplies)
	rg.GET("/movies/:id/comments", commentHandler.GetMovieComments)
	rg.GET("/cinemas/:id/comments", commentHandler.GetCinemaComments)

	rg.POST("/comments", commentHandler.CreateComment)
	rg.PUT("/comments/:id", commentHandler.UpdateComment)
	rg.DELETE("/comments/:id", commentHandler.DeleteComment)
	rg.POST("/comments/:id/vote", commentHandler.VoteComment)
	rg.DELETE("/comments/:id/vote", commentHandler.RemoveVote)
	rg.PATCH("/comments/:id/status", commentHandler.UpdateCommentStatus)

	admin := rg.Group("/admin")
	{
		admin.GET("/comments", commentHandler.GetCommentsAdmin)
		admin.GET("/movies/:id/comments", commentHandler.GetMovieCommentsAdmin)
		admin.GET("/cinemas/:id/comments", commentHandler.GetCinemaCommentsAdmin)
		admin.GET("/comments/:id/replies", commentHandler.GetCommentRepliesAdmin)
	}
}

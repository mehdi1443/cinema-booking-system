package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	service_contract "github.com/mehdi1443/cinema-booking-system/internal/application/adapter"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/cache"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/middleware"
	"github.com/mehdi1443/cinema-booking-system/internal/presentation/v1/handler"
	"github.com/mehdi1443/cinema-booking-system/pkg/logger"
)

func SetUpCommentRouter(
	rg *gin.RouterGroup,
	commentHandler *handler.CommentHandler,
	authMiddleware gin.HandlerFunc,
	rbacRepo DB.RBACInterface,
	userPermissionCacheRepo cache.UserPermissionCacheInterface,
	contexTokenService service_contract.ContextTokenServiceInterface,
	userRepo DB.UserInterface,
	userCacheConfig *bootstrap.UserPermissionCache,
	appLogger logger.Logger,
) {
	// Public Routes
	rg.GET("/comments/:id", commentHandler.GetCommentByID)
	rg.GET("/comments/:id/replies", commentHandler.GetCommentReplies)
	rg.GET("/movies/:id/comments", commentHandler.GetMovieComments)
	rg.GET("/cinemas/:id/comments", commentHandler.GetCinemaComments)

	// Protected Routes
	protected := rg.Group("")
	protected.Use(authMiddleware)
	{
		protected.POST("/comments",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:create", appLogger),
			commentHandler.CreateComment,
		)

		protected.PUT("/comments/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:update", appLogger),
			commentHandler.UpdateComment,
		)

		protected.DELETE("/comments/:id",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:delete", appLogger),
			commentHandler.DeleteComment,
		)

		protected.POST("/comments/:id/vote",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:vote", appLogger),
			commentHandler.VoteComment,
		)

		protected.DELETE("/comments/:id/vote",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:remove_vote", appLogger),
			commentHandler.RemoveVote,
		)

		protected.PATCH("/comments/:id/status",
			middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:manage_status", appLogger),
			commentHandler.UpdateCommentStatus,
		)

		// Admin / Support Routes
		admin := protected.Group("/admin")
		{
			admin.GET("/comments",
				middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:admin_read", appLogger),
				commentHandler.GetCommentsAdmin,
			)

			admin.GET("/movies/:id/comments",
				middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:admin_read", appLogger),
				commentHandler.GetMovieCommentsAdmin,
			)

			admin.GET("/cinemas/:id/comments",
				middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:admin_read", appLogger),
				commentHandler.GetCinemaCommentsAdmin,
			)

			admin.GET("/comments/:id/replies",
				middleware.ControlUserAccess(rbacRepo, userPermissionCacheRepo, contexTokenService, userRepo, userCacheConfig, "comments:admin_read", appLogger),
				commentHandler.GetCommentRepliesAdmin,
			)
		}
	}
}

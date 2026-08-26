package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type commentRepository struct {
	db database.Database
}

func NewCommentRepository(db database.Database) DB.CommentInterface {
	return &commentRepository{db: db}
}

func handleCommentDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("duplicate comment key: %w", exception.ErrCommentAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("postgres unique violation: %w", exception.ErrCommentAlreadyExists)
		case "23503": // foreign_key_violation
			constraint := strings.ToLower(pgErr.ConstraintName)
			detail := strings.ToLower(pgErr.Detail)

			if strings.Contains(constraint, "user") || strings.Contains(detail, "user_id") {
				return fmt.Errorf("user foreign key violation: %w", exception.ErrUserNotFound)
			}
			if strings.Contains(constraint, "movie") || strings.Contains(detail, "movie_id") {
				return fmt.Errorf("movie foreign key violation: %w", exception.ErrMovieNotFound)
			}
			if strings.Contains(constraint, "cinema") || strings.Contains(detail, "cinema_id") {
				return fmt.Errorf("cinema foreign key violation: %w", exception.ErrCinemaNotFound)
			}
			if strings.Contains(constraint, "parent") || strings.Contains(detail, "parent_id") {
				return fmt.Errorf("parent comment foreign key violation: %w", exception.ErrCommentNotFound)
			}
			return fmt.Errorf("unknown foreign key violation: %w", exception.ErrDatabase)
		case "23502", "23514", "40001", "40P01":
			return fmt.Errorf("database constraint/transaction error: %w", exception.ErrDatabase)
		}
	}

	return fmt.Errorf("comment database operation failed: %w", exception.ErrDatabase)
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	err := r.db.GetDB().WithContext(ctx).Omit("ID", "User", "Movie", "Cinema", "Parent", "Replies").Create(comment).Error
	if err != nil {
		return nil, handleCommentDBError(err)
	}

	return r.GetByID(ctx, comment.ID)
}

func (r *commentRepository) GetByID(ctx context.Context, id int64) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Replies.User").
		First(&comment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrCommentNotFound
		}
		return nil, handleCommentDBError(err)
	}

	return &comment, nil
}

func (r *commentRepository) GetByMovieID(ctx context.Context, movieID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error) {
	var comments []entity.Comment
	query := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("movie_id = ? AND parent_id IS NULL", movieID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
		query = query.Preload("Replies", "status IN ?", statuses).Preload("Replies.User")
	} else {
		query = query.Preload("Replies.User")
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, handleCommentDBError(err)
	}

	return comments, nil
}

func (r *commentRepository) GetByCinemaID(ctx context.Context, cinemaID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error) {
	var comments []entity.Comment
	query := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("cinema_id = ? AND parent_id IS NULL", cinemaID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
		query = query.Preload("Replies", "status IN ?", statuses).Preload("Replies.User")
	} else {
		query = query.Preload("Replies.User")
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, handleCommentDBError(err)
	}

	return comments, nil
}

func (r *commentRepository) GetRepliesByParentID(ctx context.Context, parentID int64, limit, offset int, statuses ...enum.CommentStatus) ([]entity.Comment, error) {
	var replies []entity.Comment
	query := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("parent_id = ?", parentID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	err := query.Order("created_at ASC").Limit(limit).Offset(offset).Find(&replies).Error
	if err != nil {
		return nil, handleCommentDBError(err)
	}

	return replies, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	updateData := map[string]interface{}{
		"body":       comment.Body,
		"status":     comment.Status,
		"updated_at": comment.UpdatedAt,
	}

	res := r.db.GetDB().WithContext(ctx).
		Model(&entity.Comment{}).
		Where("id = ?", comment.ID).
		Updates(updateData)

	if res.Error != nil {
		return nil, handleCommentDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrCommentNotFound
	}

	return r.GetByID(ctx, comment.ID)
}

func (r *commentRepository) UpdateStatus(ctx context.Context, id int64, status enum.CommentStatus) error {
	res := r.db.GetDB().WithContext(ctx).
		Model(&entity.Comment{}).
		Where("id = ?", id).
		Update("status", status)

	if res.Error != nil {
		return handleCommentDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrCommentNotFound
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Comment{}, id)
	if res.Error != nil {
		return handleCommentDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrCommentNotFound
	}

	return nil
}

func (r *commentRepository) AddVote(ctx context.Context, vote *entity.CommentVote) error {
	return r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		var existing entity.CommentVote
		err := tx.Where("user_id = ? AND comment_id = ?", vote.UserID, vote.CommentID).First(&existing).Error

		if err == nil {
			if existing.Vote == vote.Vote {
				return nil
			}

			if err := tx.Model(&existing).Where("user_id = ? AND comment_id = ?", vote.UserID, vote.CommentID).Update("vote", vote.Vote).Error; err != nil {
				return handleCommentDBError(err)
			}

			var likesDiff, dislikesDiff int
			if vote.Vote == enum.VoteTypeLike {
				likesDiff, dislikesDiff = 1, -1
			} else {
				likesDiff, dislikesDiff = -1, 1
			}

			err = tx.Model(&entity.Comment{}).Where("id = ?", vote.CommentID).Updates(map[string]interface{}{
				"likes_count":    gorm.Expr("likes_count + ?", likesDiff),
				"dislikes_count": gorm.Expr("dislikes_count + ?", dislikesDiff),
			}).Error
			if err != nil {
				return handleCommentDBError(err)
			}

			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return handleCommentDBError(err)
		}

		if err := tx.Create(vote).Error; err != nil {
			return handleCommentDBError(err)
		}

		column := "likes_count"
		if vote.Vote == enum.VoteTypeDislike {
			column = "dislikes_count"
		}

		err = tx.Model(&entity.Comment{}).Where("id = ?", vote.CommentID).
			UpdateColumn(column, gorm.Expr(column+" + ?", 1)).Error
		if err != nil {
			return handleCommentDBError(err)
		}

		return nil
	})
}

func (r *commentRepository) RemoveVote(ctx context.Context, userID int64, commentID int64) error {
	return r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		var existing entity.CommentVote
		err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return handleCommentDBError(err)
		}

		if err := tx.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&entity.CommentVote{}).Error; err != nil {
			return handleCommentDBError(err)
		}

		column := "likes_count"
		if existing.Vote == enum.VoteTypeDislike {
			column = "dislikes_count"
		}

		err = tx.Model(&entity.Comment{}).Where("id = ?", commentID).
			UpdateColumn(column, gorm.Expr("GREATEST("+column+" - 1, 0)")).Error
		if err != nil {
			return handleCommentDBError(err)
		}

		return nil
	})
}

func (r *commentRepository) GetUserVote(ctx context.Context, userID int64, commentID int64) (*entity.CommentVote, error) {
	var vote entity.CommentVote
	err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&vote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, handleCommentDBError(err)
	}

	return &vote, nil
}

func (r *commentRepository) GetListAdmin(ctx context.Context, status *enum.CommentStatus, limit, offset int) ([]entity.Comment, error) {
	var comments []entity.Comment
	query := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Movie").
		Preload("Cinema")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, handleCommentDBError(err)
	}

	return comments, nil
}

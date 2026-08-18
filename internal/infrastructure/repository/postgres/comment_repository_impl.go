package postgres_repo

import (
	"context"
	"errors"
	"fmt"

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

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	err := r.db.GetDB().WithContext(ctx).Omit("ID", "User", "Movie", "Cinema", "Parent", "Replies").Create(comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return &comment, nil
}

func (r *commentRepository) GetByMovieID(ctx context.Context, movieID int64, limit, offset int) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Replies.User").
		Where("movie_id = ? AND parent_id IS NULL", movieID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return comments, nil
}

func (r *commentRepository) GetByCinemaID(ctx context.Context, cinemaID int64, limit, offset int) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Preload("Replies.User").
		Where("cinema_id = ? AND parent_id IS NULL", cinemaID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return comments, nil
}

func (r *commentRepository) GetRepliesByParentID(ctx context.Context, parentID int64, limit, offset int) ([]entity.Comment, error) {
	var replies []entity.Comment
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
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
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrCommentNotFound
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Comment{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
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

			if err := tx.Model(&existing).Update("vote", vote.Vote).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}

			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if err := tx.Create(vote).Error; err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		column := "likes_count"
		if vote.Vote == enum.VoteTypeDislike {
			column = "dislikes_count"
		}

		err = tx.Model(&entity.Comment{}).Where("id = ?", vote.CommentID).
			UpdateColumn(column, gorm.Expr(column+" + ?", 1)).Error
		if err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if err := tx.Delete(&existing).Error; err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		column := "likes_count"
		if existing.Vote == enum.VoteTypeDislike {
			column = "dislikes_count"
		}

		err = tx.Model(&entity.Comment{}).Where("id = ?", commentID).
			UpdateColumn(column, gorm.Expr("GREATEST("+column+" - 1, 0)")).Error
		if err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return &vote, nil
}

package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ratingRepository struct {
	db database.Database
}

func NewRatingRepository(db database.Database) DB.RatingInterface {
	return &ratingRepository{db: db}
}

func handleRatingDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("duplicate rating key: %w", exception.ErrRatingAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("postgres unique violation: %w", exception.ErrRatingAlreadyExists)
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
			return fmt.Errorf("unknown foreign key violation: %w", exception.ErrDatabase)
		case "23502", "23514", "40001", "40P01":
			return fmt.Errorf("database constraint/transaction error: %w", exception.ErrDatabase)
		}
	}

	return fmt.Errorf("rating database operation failed: %w", exception.ErrDatabase)
}

func (r *ratingRepository) Upsert(ctx context.Context, rating *entity.Rating) (*entity.Rating, error) {
	db := r.db.GetDB().WithContext(ctx)

	var existing entity.Rating
	var err error

	if rating.MovieID != nil {
		err = db.Where("user_id = ? AND movie_id = ?", rating.UserID, *rating.MovieID).First(&existing).Error
	} else if rating.CinemaID != nil {
		err = db.Where("user_id = ? AND cinema_id = ? AND movie_id IS NULL", rating.UserID, *rating.CinemaID).First(&existing).Error
	} else {
		return nil, exception.ErrInvalidInput
	}

	if err == nil {
		existing.Score = rating.Score
		if err := db.Model(&existing).Select("score", "updated_at").Updates(existing).Error; err != nil {
			return nil, handleRatingDBError(err)
		}
		return r.GetByID(ctx, existing.ID)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Omit("User", "Movie", "Cinema").Create(rating).Error; err != nil {
			return nil, handleRatingDBError(err)
		}
		return r.GetByID(ctx, rating.ID)
	}

	return nil, handleRatingDBError(err)
}

func (r *ratingRepository) GetByID(ctx context.Context, id int64) (*entity.Rating, error) {
	var rating entity.Rating
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		First(&rating, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrRatingNotFound
		}
		return nil, handleRatingDBError(err)
	}

	return &rating, nil
}

func (r *ratingRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Rating{}, id)
	if res.Error != nil {
		return handleRatingDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrRatingNotFound
	}

	return nil
}

func (r *ratingRepository) GetUserRatingForMovie(ctx context.Context, userID, movieID int64) (*entity.Rating, error) {
	var rating entity.Rating
	err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		First(&rating).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil indicating not rated yet
		}
		return nil, handleRatingDBError(err)
	}

	return &rating, nil
}

func (r *ratingRepository) GetUserRatingForCinema(ctx context.Context, userID, cinemaID int64) (*entity.Rating, error) {
	var rating entity.Rating
	err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND cinema_id = ? AND movie_id IS NULL", userID, cinemaID).
		First(&rating).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil indicating not rated yet
		}
		return nil, handleRatingDBError(err)
	}

	return &rating, nil
}

func (r *ratingRepository) GetByMovieID(ctx context.Context, movieID int64, limit, offset int) ([]entity.Rating, error) {
	var ratings []entity.Rating
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("movie_id = ?", movieID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).Error

	if err != nil {
		return nil, handleRatingDBError(err)
	}

	return ratings, nil
}

func (r *ratingRepository) GetByCinemaID(ctx context.Context, cinemaID int64, limit, offset int) ([]entity.Rating, error) {
	var ratings []entity.Rating
	err := r.db.GetDB().WithContext(ctx).
		Preload("User").
		Where("cinema_id = ? AND movie_id IS NULL", cinemaID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).Error

	if err != nil {
		return nil, handleRatingDBError(err)
	}

	return ratings, nil
}

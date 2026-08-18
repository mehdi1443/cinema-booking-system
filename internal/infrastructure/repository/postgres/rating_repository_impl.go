package postgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ratingRepository struct {
	db database.Database
}

func NewRatingRepository(db database.Database) DB.RatingInterface {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Upsert(ctx context.Context, rating *entity.Rating) (*entity.Rating, error) {
	err := r.db.GetDB().WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "movie_id"}, {Name: "cinema_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"score", "updated_at"}),
		}).
		Omit("User", "Movie", "Cinema").
		Create(rating).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return r.GetByID(ctx, rating.ID)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return &rating, nil
}

func (r *ratingRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Rating{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
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
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
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
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return ratings, nil
}

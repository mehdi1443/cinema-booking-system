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
)

type favoriteRepository struct {
	db database.Database
}

func NewFavoriteRepository(db database.Database) DB.FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) AddMovieToFavorites(ctx context.Context, userID int64, movieID int64) (*entity.Favorite, error) {
	favorite := &entity.Favorite{
		UserID:  userID,
		MovieID: &movieID,
	}

	err := r.db.GetDB().WithContext(ctx).Create(favorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exception.ErrMovieNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return r.getByID(ctx, favorite.ID)
}

func (r *favoriteRepository) RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Delete(&entity.Favorite{})

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrRecordNotFound
	}

	return nil
}

func (r *favoriteRepository) IsMovieFavorite(ctx context.Context, userID int64, movieID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Model(&entity.Favorite{}).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return count > 0, nil
}

func (r *favoriteRepository) AddCinemaToFavorites(ctx context.Context, userID int64, cinemaID int64) (*entity.Favorite, error) {
	favorite := &entity.Favorite{
		UserID:   userID,
		CinemaID: &cinemaID,
	}

	err := r.db.GetDB().WithContext(ctx).Create(favorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exception.ErrCinemaNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return r.getByID(ctx, favorite.ID)
}

func (r *favoriteRepository) RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND cinema_id = ?", userID, cinemaID).
		Delete(&entity.Favorite{})

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrRecordNotFound
	}

	return nil
}

func (r *favoriteRepository) IsCinemaFavorite(ctx context.Context, userID int64, cinemaID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Model(&entity.Favorite{}).
		Where("user_id = ? AND cinema_id = ?", userID, cinemaID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return count > 0, nil
}

func (r *favoriteRepository) GetUserFavorites(ctx context.Context, filter DB.FavoriteFilter) ([]entity.Favorite, error) {
	var favorites []entity.Favorite
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Favorite{}).Where("user_id = ?", filter.UserID)

	query = r.applyFilters(query, filter)

	err := query.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Preload("Movie").
		Preload("Cinema").
		Find(&favorites).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return favorites, nil
}

func (r *favoriteRepository) CountUserFavorites(ctx context.Context, filter DB.FavoriteFilter) (int64, error) {
	var count int64
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Favorite{}).Where("user_id = ?", filter.UserID)

	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return count, nil
}

func (r *favoriteRepository) getByID(ctx context.Context, id int64) (*entity.Favorite, error) {
	var favorite entity.Favorite
	err := r.db.GetDB().WithContext(ctx).
		Preload("Movie").
		Preload("Cinema").
		First(&favorite, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &favorite, nil
}

func (r *favoriteRepository) applyFilters(query *gorm.DB, filter DB.FavoriteFilter) *gorm.DB {
	switch filter.Target {
	case "movie":
		query = query.Where("movie_id IS NOT NULL")
	case "cinema":
		query = query.Where("cinema_id IS NOT NULL")
	}
	return query
}

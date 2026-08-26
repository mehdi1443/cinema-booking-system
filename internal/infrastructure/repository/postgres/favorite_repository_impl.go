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

type favoriteRepository struct {
	db database.Database
}

func NewFavoriteRepository(db database.Database) DB.FavoriteRepository {
	return &favoriteRepository{db: db}
}

func handleFavoriteDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("duplicate favorite key: %w", exception.ErrFavoriteAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("postgres unique violation: %w", exception.ErrFavoriteAlreadyExists)
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

	return fmt.Errorf("favorite database operation failed: %w", exception.ErrDatabase)
}

func (r *favoriteRepository) AddMovieToFavorites(ctx context.Context, userID int64, movieID int64) (*entity.Favorite, error) {
	favorite := &entity.Favorite{
		UserID:  userID,
		MovieID: &movieID,
	}

	err := r.db.GetDB().WithContext(ctx).Create(favorite).Error
	if err != nil {
		return nil, handleFavoriteDBError(err)
	}

	return r.getByID(ctx, favorite.ID)
}

func (r *favoriteRepository) RemoveMovieFromFavorites(ctx context.Context, userID int64, movieID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Delete(&entity.Favorite{})

	if res.Error != nil {
		return handleFavoriteDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrFavoriteNotFound
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
		return false, handleFavoriteDBError(err)
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
		return nil, handleFavoriteDBError(err)
	}

	return r.getByID(ctx, favorite.ID)
}

func (r *favoriteRepository) RemoveCinemaFromFavorites(ctx context.Context, userID int64, cinemaID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND cinema_id = ?", userID, cinemaID).
		Delete(&entity.Favorite{})

	if res.Error != nil {
		return handleFavoriteDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrFavoriteNotFound
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
		return false, handleFavoriteDBError(err)
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
		return nil, handleFavoriteDBError(err)
	}

	return favorites, nil
}

func (r *favoriteRepository) CountUserFavorites(ctx context.Context, filter DB.FavoriteFilter) (int64, error) {
	var count int64
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Favorite{}).Where("user_id = ?", filter.UserID)

	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	if err != nil {
		return 0, handleFavoriteDBError(err)
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
			return nil, exception.ErrFavoriteNotFound
		}
		return nil, handleFavoriteDBError(err)
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

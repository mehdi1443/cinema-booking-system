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

type genreRepository struct {
	db database.Database
}

func NewGenreRepository(db database.Database) DB.GenreInterface {
	return &genreRepository{db: db}
}

func (r *genreRepository) Create(ctx context.Context, genre *entity.Genre) (*entity.Genre, error) {
	err := r.db.GetDB().WithContext(ctx).Create(genre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return genre, nil
}

func (r *genreRepository) GetByID(ctx context.Context, id int64) (*entity.Genre, error) {
	var genre entity.Genre
	err := r.db.GetDB().WithContext(ctx).First(&genre, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrGenreNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &genre, nil
}

func (r *genreRepository) GetByName(ctx context.Context, name string) (*entity.Genre, error) {
	var genre entity.Genre
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&genre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrGenreNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &genre, nil
}

func (r *genreRepository) GetGenres(ctx context.Context, filter DB.GenreFilter, opts *DB.QueryOptions) ([]entity.Genre, int64, error) {
	var genres []entity.Genre
	q := r.db.GetDB().WithContext(ctx).Model(&entity.Genre{})

	if filter.SearchKeyword != "" {
		q = q.Where("name ILIKE ?", "%"+filter.SearchKeyword+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	q = q.Order("name ASC")

	if opts != nil && opts.HasPagination() {
		q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	if err := q.Find(&genres).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return genres, total, nil
}

func (r *genreRepository) Update(ctx context.Context, genre *entity.Genre) (*entity.Genre, error) {
	res := r.db.GetDB().WithContext(ctx).Model(genre).Where("id = ?", genre.ID).Updates(genre)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrGenreNotFound
	}

	return genre, nil
}

func (r *genreRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Genre{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrGenreNotFound
	}

	return nil
}

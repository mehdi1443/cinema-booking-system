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

func (r *genreRepository) GetAll(ctx context.Context) ([]entity.Genre, error) {
	var genres []entity.Genre
	err := r.db.GetDB().WithContext(ctx).Find(&genres).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return genres, nil
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

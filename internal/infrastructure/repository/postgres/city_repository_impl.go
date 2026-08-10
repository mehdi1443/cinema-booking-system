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

type cityRepository struct {
	db database.Database
}

func NewCityRepository(db database.Database) DB.CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) Create(ctx context.Context, city *entity.City) error {
	err := r.db.GetDB().WithContext(ctx).Create(city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrDuplicateKey
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *cityRepository) GetByID(ctx context.Context, id int64) (*entity.City, error) {
	var city entity.City
	err := r.db.GetDB().WithContext(ctx).Preload("Province").First(&city, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrCityNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &city, nil
}

func (r *cityRepository) GetByNameAndProvinceID(ctx context.Context, name string, provinceID int64) (*entity.City, error) {
	var city entity.City
	err := r.db.GetDB().WithContext(ctx).
		Where("name = ? AND province_id = ?", name, provinceID).
		First(&city).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrCityNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &city, nil
}

func (r *cityRepository) GetByProvinceID(ctx context.Context, provinceID int64) ([]entity.City, error) {
	var cities []entity.City
	err := r.db.GetDB().WithContext(ctx).
		Where("province_id = ?", provinceID).
		Find(&cities).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return cities, nil
}

func (r *cityRepository) GetAll(ctx context.Context) ([]entity.City, error) {
	var cities []entity.City
	err := r.db.GetDB().WithContext(ctx).Preload("Province").Find(&cities).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return cities, nil
}

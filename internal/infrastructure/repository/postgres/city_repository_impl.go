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

func (r *cityRepository) Update(ctx context.Context, city *entity.City) (*entity.City, error) {
	city.Province = nil

	err := r.db.GetDB().WithContext(ctx).
		Model(&entity.City{}).
		Where("id = ?", city.ID).
		Select("name", "province_id").
		Updates(map[string]interface{}{
			"name":        city.Name,
			"province_id": city.ProvinceID,
		}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return city, nil
}

func (r *cityRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.GetDB().WithContext(ctx).Delete(&entity.City{}, id)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.ErrCityNotFound
	}
	return nil
}

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

type provinceRepository struct {
	db database.Database
}

func NewProvinceRepository(db database.Database) DB.ProvinceRepository {
	return &provinceRepository{db: db}
}

func (r *provinceRepository) Create(ctx context.Context, province *entity.Province) error {
	err := r.db.GetDB().WithContext(ctx).Create(province).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrDuplicateKey
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

func (r *provinceRepository) GetByID(ctx context.Context, id int64) (*entity.Province, error) {
	var province entity.Province
	err := r.db.GetDB().WithContext(ctx).First(&province, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrProvinceNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &province, nil
}

func (r *provinceRepository) GetByName(ctx context.Context, name string) (*entity.Province, error) {
	var province entity.Province
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&province).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrProvinceNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &province, nil
}

func (r *provinceRepository) GetAll(ctx context.Context) ([]entity.Province, error) {
	var provinces []entity.Province
	err := r.db.GetDB().WithContext(ctx).Find(&provinces).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return provinces, nil
}

func (r *provinceRepository) Update(ctx context.Context, province *entity.Province) (*entity.Province, error) {
	err := r.db.GetDB().WithContext(ctx).
		Model(&entity.Province{}).
		Where("id = ?", province.ID).
		Select("name").
		Updates(map[string]interface{}{
			"name": province.Name,
		}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return province, nil
}

func (r *provinceRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.GetDB().WithContext(ctx).Delete(&entity.Province{}, id)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.ErrProvinceNotFound
	}
	return nil
}

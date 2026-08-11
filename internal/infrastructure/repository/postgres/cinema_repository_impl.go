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

type CinemaFacility struct {
	CinemaID   int64 `gorm:"primaryKey;column:cinema_id"`
	FacilityID int64 `gorm:"primaryKey;column:facility_id"`
}

func (CinemaFacility) TableName() string {
	return "cinema_facilities"
}

type cinemaRepository struct {
	db database.Database
}

func NewCinemaRepository(db database.Database) DB.CinemaRepository {
	return &cinemaRepository{db: db}
}

func (r *cinemaRepository) Create(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		if err := tx.Omit("ID", "City", "Facilities").Create(cinema).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return exception.ErrCityNotFound
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if len(cinema.Facilities) > 0 {
			var cinemaFacilities []CinemaFacility
			for _, facility := range cinema.Facilities {
				cinemaFacilities = append(cinemaFacilities, CinemaFacility{
					CinemaID:   cinema.ID,
					FacilityID: facility.ID,
				})
			}
			if err := tx.Create(&cinemaFacilities).Error; err != nil {
				if errors.Is(err, gorm.ErrForeignKeyViolated) {
					return exception.ErrFacilityNotFound
				}
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, cinema.ID)
}

func (r *cinemaRepository) GetByID(ctx context.Context, id int64) (*entity.Cinema, error) {
	var cinema entity.Cinema
	err := r.db.GetDB().WithContext(ctx).
		Preload("City").
		Preload("City.Province").
		Preload("Facilities").
		First(&cinema, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrCinemaNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &cinema, nil
}

func (r *cinemaRepository) Update(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		updateData := map[string]interface{}{
			"city_id":        cinema.CityID,
			"name":           cinema.Name,
			"description":    cinema.Description,
			"address":        cinema.Address,
			"latitude":       cinema.Latitude,
			"longitude":      cinema.Longitude,
			"phone":          cinema.Phone,
			"image_url_low":  cinema.ImageURLLow,
			"image_url_high": cinema.ImageURLHigh,
			"status":         cinema.Status,
			"updated_at":     cinema.UpdatedAt,
		}

		res := tx.Model(&entity.Cinema{}).Where("id = ?", cinema.ID).Updates(updateData)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			if errors.Is(res.Error, gorm.ErrForeignKeyViolated) {
				return exception.ErrCityNotFound
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrCinemaNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, cinema.ID)
}

func (r *cinemaRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		if err := tx.Where("cinema_id = ?", id).Delete(&CinemaFacility{}).Error; err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		res := tx.Delete(&entity.Cinema{}, id)
		if res.Error != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrCinemaNotFound
		}

		return nil
	})
}

func (r *cinemaRepository) GetWithFilters(ctx context.Context, filter DB.CinemaFilter) ([]entity.Cinema, error) {
	var cinemas []entity.Cinema

	subQuery := r.db.GetDB().WithContext(ctx).Table("cinemas").Select("cinemas.id")
	subQuery = r.applyFilters(subQuery, filter)

	err := r.db.GetDB().WithContext(ctx).
		Where("id IN (?)", subQuery).
		Limit(filter.Limit).
		Offset(filter.Offset).
		Preload("City").
		Preload("City.Province").
		Preload("Facilities").
		Find(&cinemas).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return cinemas, nil
}

func (r *cinemaRepository) CountWithFilters(ctx context.Context, filter DB.CinemaFilter) (int64, error) {
	var count int64
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Cinema{})

	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return count, nil
}

func (r *cinemaRepository) AddFacility(ctx context.Context, cinemaID int64, facilityID int64) error {
	cinemaFacility := CinemaFacility{
		CinemaID:   cinemaID,
		FacilityID: facilityID,
	}

	err := r.db.GetDB().WithContext(ctx).Create(&cinemaFacility).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return exception.ErrDuplicateKey
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return exception.ErrFacilityNotFound
		}
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return nil
}

func (r *cinemaRepository) RemoveFacility(ctx context.Context, cinemaID int64, facilityID int64) error {
	res := r.db.GetDB().WithContext(ctx).
		Where("cinema_id = ? AND facility_id = ?", cinemaID, facilityID).
		Delete(&CinemaFacility{})

	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrFacilityNotFound
	}

	return nil
}

func (r *cinemaRepository) UpdateFacilities(ctx context.Context, cinemaID int64, facilityIDs []int64) error {
	return r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		if err := tx.Where("cinema_id = ?", cinemaID).Delete(&CinemaFacility{}).Error; err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if len(facilityIDs) > 0 {
			var cinemaFacilities []CinemaFacility
			for _, fID := range facilityIDs {
				cinemaFacilities = append(cinemaFacilities, CinemaFacility{
					CinemaID:   cinemaID,
					FacilityID: fID,
				})
			}
			if err := tx.Create(&cinemaFacilities).Error; err != nil {
				if errors.Is(err, gorm.ErrForeignKeyViolated) {
					return exception.ErrFacilityNotFound
				}
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		return nil
	})
}

func (r *cinemaRepository) applyFilters(query *gorm.DB, filter DB.CinemaFilter) *gorm.DB {
	if filter.CityID != nil {
		query = query.Where("cinemas.city_id = ?", *filter.CityID)
	}

	if filter.ProvinceID != nil {
		query = query.Joins("JOIN cities ON cities.id = cinemas.city_id").
			Where("cities.province_id = ?", *filter.ProvinceID)
	}

	if filter.Status != nil {
		query = query.Where("cinemas.status = ?", *filter.Status)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("cinemas.name ILIKE ? OR cinemas.address ILIKE ?", searchTerm, searchTerm)
	}

	if len(filter.FacilityIDs) > 0 {
		query = query.Joins("JOIN cinema_facilities ON cinema_facilities.cinema_id = cinemas.id").
			Where("cinema_facilities.facility_id IN ?", filter.FacilityIDs).
			Group("cinemas.id").
			Having("COUNT(DISTINCT cinema_facilities.facility_id) = ?", len(filter.FacilityIDs))
	}

	return query
}

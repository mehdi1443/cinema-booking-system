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

func NewCinemaRepository(db database.Database) DB.CinemaInterface {
	return &cinemaRepository{db: db}
}

func handleCinemaDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return exception.ErrDuplicateKey
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return exception.ErrDuplicateKey
		case "23503": // foreign_key_violation
			if strings.Contains(pgErr.ConstraintName, "city") {
				return exception.ErrCityNotFound
			}
			if strings.Contains(pgErr.ConstraintName, "facility") || strings.Contains(pgErr.Detail, "cinema_facilities") {
				return exception.ErrFacilityNotFound
			}
			return exception.ErrDatabase // Fallback
		}
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return exception.ErrDatabase
	}

	return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
}

func (r *cinemaRepository) Create(ctx context.Context, cinema *entity.Cinema, facilityIDs []int64) (*entity.Cinema, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		if err := tx.Omit("ID", "City", "Facilities").Create(cinema).Error; err != nil {
			return handleCinemaDBError(err)
		}

		if len(facilityIDs) > 0 {
			var cinemaFacilities []CinemaFacility
			for _, fID := range facilityIDs {
				cinemaFacilities = append(cinemaFacilities, CinemaFacility{
					CinemaID:   cinema.ID,
					FacilityID: fID,
				})
			}
			if err := tx.Create(&cinemaFacilities).Error; err != nil {
				return handleCinemaDBError(err)
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
		return nil, handleCinemaDBError(err)
	}
	return &cinema, nil
}

func (r *cinemaRepository) GetCinemas(ctx context.Context, filter DB.CinemaFilter, opts *DB.QueryOptions) ([]entity.Cinema, int64, error) {
	var cinemas []entity.Cinema

	db := r.db.GetDB().WithContext(ctx)
	q := db.Model(&entity.Cinema{})

	if filter.SearchKeyword != "" {
		searchTerm := "%" + filter.SearchKeyword + "%"
		q = q.Where("cinemas.name ILIKE ? OR cinemas.address ILIKE ?", searchTerm, searchTerm)
	}

	if filter.CityID != nil {
		q = q.Where("cinemas.city_id = ?", *filter.CityID)
	}

	if filter.ProvinceID != nil {
		q = q.Joins("JOIN cities ON cities.id = cinemas.city_id").
			Where("cities.province_id = ?", *filter.ProvinceID)
	}

	if filter.Status != nil {
		q = q.Where("cinemas.status = ?", *filter.Status)
	}

	if filter.MinRating != nil {
		q = q.Where("cinemas.rating_avg >= ?", *filter.MinRating)
	}

	if filter.MaxRating != nil {
		q = q.Where("cinemas.rating_avg <= ?", *filter.MaxRating)
	}

	if len(filter.FacilityIDs) > 0 {
		subQuery := db.Table("cinema_facilities").
			Select("cinema_id").
			Where("facility_id IN ?", filter.FacilityIDs).
			Group("cinema_id").
			Having("COUNT(DISTINCT facility_id) = ?", len(filter.FacilityIDs))

		q = q.Where("cinemas.id IN (?)", subQuery)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, handleCinemaDBError(err)
	}

	q = q.Preload("City").Preload("City.Province").Preload("Facilities")

	if opts != nil && opts.HasSorting() {
		direction := "ASC"
		if !opts.Sorting.Asc {
			direction = "DESC"
		}
		q = q.Order(fmt.Sprintf("cinemas.%s %s", opts.Sorting.Column, direction))
	} else {
		q = q.Order("cinemas.created_at DESC")
	}

	if opts != nil && opts.HasPagination() {
		q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	if err := q.Find(&cinemas).Error; err != nil {
		return nil, 0, handleCinemaDBError(err)
	}

	return cinemas, total, nil
}

func (r *cinemaRepository) Update(ctx context.Context, cinema *entity.Cinema, facilityIDs []int64) (*entity.Cinema, error) {
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
			return handleCinemaDBError(res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrCinemaNotFound
		}

		if facilityIDs != nil {
			if err := tx.Where("cinema_id = ?", cinema.ID).Delete(&CinemaFacility{}).Error; err != nil {
				return handleCinemaDBError(err)
			}

			if len(facilityIDs) > 0 {
				var cinemaFacilities []CinemaFacility
				for _, fID := range facilityIDs {
					cinemaFacilities = append(cinemaFacilities, CinemaFacility{
						CinemaID:   cinema.ID,
						FacilityID: fID,
					})
				}
				if err := tx.Create(&cinemaFacilities).Error; err != nil {
					return handleCinemaDBError(err)
				}
			}
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
			return handleCinemaDBError(err)
		}

		res := tx.Delete(&entity.Cinema{}, id)
		if res.Error != nil {
			return handleCinemaDBError(res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrCinemaNotFound
		}

		return nil
	})
}

func (r *cinemaRepository) UpdateRatingAggregate(ctx context.Context, cinemaID int64) error {
	query := `
        UPDATE cinemas
        SET 
            rating_avg = COALESCE((SELECT ROUND(AVG(score)::numeric, 2) FROM ratings WHERE movie_id IS NULL AND cinema_id = ?), 0.00),
            rating_count = (SELECT COUNT(*) FROM ratings WHERE movie_id IS NULL AND cinema_id = ?)
        WHERE id = ?
    `
	err := r.db.GetDB().WithContext(ctx).Exec(query, cinemaID, cinemaID, cinemaID).Error
	if err != nil {
		return handleCinemaDBError(err)
	}
	return nil
}

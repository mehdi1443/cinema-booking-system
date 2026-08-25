package postgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type facilityRepository struct {
	db database.Database
}

func NewFacilityRepository(db database.Database) DB.FacilityRepository {
	return &facilityRepository{db: db}
}

func handleFacilityDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("duplicate facility key: %w", exception.ErrFacilityAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("postgres unique violation: %w", exception.ErrFacilityAlreadyExists)
		case "23503": // foreign_key_violation
			return fmt.Errorf("facility is in use: %w", exception.ErrResourceInUse)
		case "23502": // not_null_violation
		case "23514": // check_violation
		case "40001": // serialization_failure
		case "40P01": // deadlock_detected
			return fmt.Errorf("database constraint/transaction error: %w", exception.ErrDatabase)
		}
	}

	return fmt.Errorf("facility database operation failed: %w", exception.ErrDatabase)
}

func (r *facilityRepository) Create(ctx context.Context, facility *entity.Facility) (*entity.Facility, error) {
	err := r.db.GetDB().WithContext(ctx).Create(facility).Error
	if err != nil {
		return nil, handleFacilityDBError(err)
	}
	return facility, nil
}

func (r *facilityRepository) GetByID(ctx context.Context, id int64) (*entity.Facility, error) {
	var facility entity.Facility
	err := r.db.GetDB().WithContext(ctx).First(&facility, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrFacilityNotFound
		}
		return nil, handleFacilityDBError(err)
	}
	return &facility, nil
}

func (r *facilityRepository) GetAll(ctx context.Context) ([]entity.Facility, error) {
	var facilities []entity.Facility
	err := r.db.GetDB().WithContext(ctx).Find(&facilities).Error
	if err != nil {
		return nil, handleFacilityDBError(err)
	}
	return facilities, nil
}

func (r *facilityRepository) GetByCinemaID(ctx context.Context, cinemaID int64) ([]entity.Facility, error) {
	var facilities []entity.Facility
	err := r.db.GetDB().WithContext(ctx).
		Joins("JOIN cinema_facilities ON cinema_facilities.facility_id = facilities.id").
		Where("cinema_facilities.cinema_id = ?", cinemaID).
		Find(&facilities).Error
	if err != nil {
		return nil, handleFacilityDBError(err)
	}
	return facilities, nil
}

func (r *facilityRepository) Update(ctx context.Context, facility *entity.Facility) (*entity.Facility, error) {
	res := r.db.GetDB().WithContext(ctx).Model(facility).Where("id = ?", facility.ID).Updates(facility)
	if res.Error != nil {
		return nil, handleFacilityDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrFacilityNotFound
	}

	return facility, nil
}

func (r *facilityRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Facility{}, id)
	if res.Error != nil {
		return handleFacilityDBError(res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrFacilityNotFound
	}

	return nil
}

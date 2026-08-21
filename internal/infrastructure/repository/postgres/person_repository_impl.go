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

type personRepository struct {
	db database.Database
}

func NewPersonRepository(db database.Database) DB.PersonInterface {
	return &personRepository{db: db}
}

func (r *personRepository) Create(ctx context.Context, person *entity.Person) (*entity.Person, error) {
	err := r.db.GetDB().WithContext(ctx).Create(person).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return person, nil
}

func (r *personRepository) GetByID(ctx context.Context, id int64) (*entity.Person, error) {
	var person entity.Person
	err := r.db.GetDB().WithContext(ctx).First(&person, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrPersonNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &person, nil
}

func (r *personRepository) GetByName(ctx context.Context, name string) (*entity.Person, error) {
	var person entity.Person
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&person).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrPersonNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &person, nil
}

func (r *personRepository) GetPersons(ctx context.Context, filter DB.PersonFilter, opts *DB.QueryOptions) ([]entity.Person, int64, error) {
	var persons []entity.Person
	q := r.db.GetDB().WithContext(ctx).Model(&entity.Person{})

	if filter.SearchKeyword != "" {
		q = q.Where("name ILIKE ?", "%"+filter.SearchKeyword+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	q = q.Order("id DESC")

	if opts != nil && opts.HasPagination() {
		q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	if err := q.Find(&persons).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return persons, total, nil
}

func (r *personRepository) Update(ctx context.Context, person *entity.Person) (*entity.Person, error) {
	updateData := map[string]interface{}{
		"name":           person.Name,
		"biography":      person.Biography,
		"image_url_low":  person.ImageURLLow,
		"image_url_high": person.ImageURLHigh,
	}

	res := r.db.GetDB().WithContext(ctx).Model(&entity.Person{}).Where("id = ?", person.ID).Updates(updateData)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, exception.ErrDuplicateKey
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrPersonNotFound
	}

	return person, nil
}

func (r *personRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Person{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrPersonNotFound
	}

	return nil
}

package postgres_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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

func handlePersonDBError(err error, operation string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrPersonNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return exception.ErrPersonAlreadyExists
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return exception.ErrPersonAlreadyExists
		case "23503": // foreign_key_violation
			return exception.ErrResourceInUse
		}
	}

	return exception.Wrap(exception.ErrDatabase, "database error during "+operation, err)
}

func (r *personRepository) Create(ctx context.Context, person *entity.Person) (*entity.Person, error) {
	err := r.db.GetDB().WithContext(ctx).Create(person).Error
	if err != nil {
		return nil, handlePersonDBError(err, "CreatePerson")
	}
	return person, nil
}

func (r *personRepository) GetByID(ctx context.Context, id int64) (*entity.Person, error) {
	var person entity.Person
	err := r.db.GetDB().WithContext(ctx).First(&person, id).Error
	if err != nil {
		return nil, handlePersonDBError(err, "GetPersonByID")
	}
	return &person, nil
}

func (r *personRepository) GetByName(ctx context.Context, name string) (*entity.Person, error) {
	var person entity.Person
	err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&person).Error
	if err != nil {
		return nil, handlePersonDBError(err, "GetPersonByName")
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
		return nil, 0, handlePersonDBError(err, "CountPersons")
	}

	q = q.Order("id DESC")

	if opts != nil && opts.HasPagination() {
		q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	if err := q.Find(&persons).Error; err != nil {
		return nil, 0, handlePersonDBError(err, "FindPersons")
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
		return nil, handlePersonDBError(res.Error, "UpdatePerson")
	}

	if res.RowsAffected == 0 {
		return nil, exception.ErrPersonNotFound
	}

	return person, nil
}

func (r *personRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Person{}, id)
	if res.Error != nil {
		return handlePersonDBError(res.Error, "DeletePerson")
	}

	if res.RowsAffected == 0 {
		return exception.ErrPersonNotFound
	}

	return nil
}

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

func (r *personRepository) GetAll(ctx context.Context) ([]entity.Person, error) {
	var people []entity.Person
	err := r.db.GetDB().WithContext(ctx).Find(&people).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return people, nil
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

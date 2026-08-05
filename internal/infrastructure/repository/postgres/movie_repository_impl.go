package postgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type MoviePerson struct {
	MovieID  int64                `gorm:"primaryKey;column:movie_id"`
	PersonID int64                `gorm:"primaryKey;column:person_id"`
	Role     enum.MoviePersonRole `gorm:"primaryKey;column:role"`
}

func (MoviePerson) TableName() string {
	return "movie_people"
}

type movieRepository struct {
	db database.Database
}

func NewMovieRepository(db database.Database) DB.MovieInterface {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		if err := tx.Create(movie).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if len(genreIDs) > 0 {
			var genres []entity.Genre
			for _, id := range genreIDs {
				genres = append(genres, entity.Genre{ID: id})
			}
			if err := tx.Model(movie).Association("Genres").Replace(genres); err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		if len(peopleIDs) > 0 {
			var moviePeople []MoviePerson
			for personID, role := range peopleIDs {
				moviePeople = append(moviePeople, MoviePerson{
					MovieID:  movie.ID,
					PersonID: personID,
					Role:     role,
				})
			}
			if err := tx.Create(&moviePeople).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, movie.ID)
}

func (r *movieRepository) GetByID(ctx context.Context, id int64) (*entity.Movie, error) {
	var movie entity.Movie
	err := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People").
		First(&movie, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrMovieNotFound
		}
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return &movie, nil
}

func (r *movieRepository) GetByStatus(ctx context.Context, status enum.MovieStatus) ([]entity.Movie, error) {
	var movies []entity.Movie
	err := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People").
		Where("status = ?", status).
		Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return movies, nil
}

func (r *movieRepository) GetByGenreID(ctx context.Context, genreID int64) ([]entity.Movie, error) {
	var movies []entity.Movie
	err := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People").
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Where("movie_genres.genre_id = ?", genreID).
		Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return movies, nil
}

func (r *movieRepository) GetAll(ctx context.Context) ([]entity.Movie, error) {
	var movies []entity.Movie
	err := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People").
		Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return movies, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		res := tx.Model(movie).Where("id = ?", movie.ID).Updates(movie)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrMovieNotFound
		}

		var genres []entity.Genre
		for _, id := range genreIDs {
			genres = append(genres, entity.Genre{ID: id})
		}
		if err := tx.Model(movie).Association("Genres").Replace(genres); err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if err := tx.Where("movie_id = ?", movie.ID).Delete(&MoviePerson{}).Error; err != nil {
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if len(peopleIDs) > 0 {
			var moviePeople []MoviePerson
			for personID, role := range peopleIDs {
				moviePeople = append(moviePeople, MoviePerson{
					MovieID:  movie.ID,
					PersonID: personID,
					Role:     role,
				})
			}
			if err := tx.Create(&moviePeople).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, movie.ID)
}

func (r *movieRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.GetDB().WithContext(ctx).Delete(&entity.Movie{}, id)
	if res.Error != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
	}

	if res.RowsAffected == 0 {
		return exception.ErrMovieNotFound
	}

	return nil
}

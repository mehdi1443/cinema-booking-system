package postgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/application/dto"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
)

type MovieGenre struct {
	MovieID int64 `gorm:"primaryKey;column:movie_id"`
	GenreID int64 `gorm:"primaryKey;column:genre_id"`
}

func (MovieGenre) TableName() string {
	return "movie_genres"
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

		if err := tx.Omit("ID", "Genres", "People").Create(movie).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
		}

		if len(genreIDs) > 0 {
			var movieGenres []MovieGenre
			for _, id := range genreIDs {
				movieGenres = append(movieGenres, MovieGenre{
					MovieID: movie.ID,
					GenreID: id,
				})
			}
			if err := tx.Create(&movieGenres).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}
		}

		if len(peopleIDs) > 0 {
			var moviePeople []entity.MoviePerson
			for personID, role := range peopleIDs {
				moviePeople = append(moviePeople, entity.MoviePerson{
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
		Preload("People.Person").
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
		Preload("People.Person").
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
		Preload("People.Person").
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Where("movie_genres.genre_id = ?", genreID).
		Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return movies, nil
}

func (r *movieRepository) GetFiltered(ctx context.Context, filter dto.MovieFilterRequest) ([]entity.Movie, error) {
	var movies []entity.Movie
	query := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People.Person")

	if filter.Status != nil {
		query = query.Where("movies.status = ?", *filter.Status)
	}

	if filter.Country != nil && *filter.Country != "" {
		query = query.Where("movies.country = ?", *filter.Country)
	}

	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("movies.title ILIKE ?", "%"+*filter.Search+"%")
	}

	if filter.GenreID != nil {
		query = query.Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
			Where("movie_genres.genre_id = ?", *filter.GenreID)
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&movies).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return movies, nil
}

func (r *movieRepository) GetAll(ctx context.Context) ([]entity.Movie, error) {
	var movies []entity.Movie
	err := r.db.GetDB().WithContext(ctx).
		Preload("Genres").
		Preload("People.Person").
		Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return movies, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *entity.Movie, genreIDs []int64, peopleIDs map[int64]enum.MoviePersonRole) (*entity.Movie, error) {
	err := r.db.WithTransaction(func(txWrapper database.Database) error {
		tx := txWrapper.GetDB().WithContext(ctx)

		updateData := map[string]interface{}{
			"title":            movie.Title,
			"description":      movie.Description,
			"duration_minutes": movie.DurationMinutes,
			"release_date":     movie.ReleaseDate,
			"country":          movie.Country,
			"poster_url_low":   movie.PosterURLLow,
			"poster_url_high":  movie.PosterURLHigh,
			"trailer_url":      movie.TrailerURL,
			"status":           movie.Status,
			"updated_at":       movie.UpdatedAt,
		}

		res := tx.Model(&entity.Movie{}).Where("id = ?", movie.ID).Updates(updateData)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
				return exception.ErrDuplicateKey
			}
			return fmt.Errorf("%w: %v", exception.ErrDatabase, res.Error)
		}
		if res.RowsAffected == 0 {
			return exception.ErrMovieNotFound
		}

		if genreIDs != nil {
			if err := tx.Where("movie_id = ?", movie.ID).Delete(&MovieGenre{}).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}

			if len(genreIDs) > 0 {
				var movieGenres []MovieGenre
				for _, id := range genreIDs {
					movieGenres = append(movieGenres, MovieGenre{
						MovieID: movie.ID,
						GenreID: id,
					})
				}
				if err := tx.Create(&movieGenres).Error; err != nil {
					return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
				}
			}
		}

		if peopleIDs != nil {
			if err := tx.Where("movie_id = ?", movie.ID).Delete(&entity.MoviePerson{}).Error; err != nil {
				return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
			}

			if len(peopleIDs) > 0 {
				var moviePeople []entity.MoviePerson
				for personID, role := range peopleIDs {
					moviePeople = append(moviePeople, entity.MoviePerson{
						MovieID:  movie.ID,
						PersonID: personID,
						Role:     role,
					})
				}
				if err := tx.Create(&moviePeople).Error; err != nil {
					return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
				}
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

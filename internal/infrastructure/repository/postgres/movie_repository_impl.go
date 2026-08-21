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

func NewMovieRepository(db database.Database) DB.MovieInterFace {
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

func (r *movieRepository) GetMovies(ctx context.Context, filter DB.MovieFilter, opts *DB.QueryOptions) ([]entity.Movie, int64, error) {
	var movies []entity.Movie

	q := r.db.GetDB().WithContext(ctx).Model(&entity.Movie{}).
		Preload("Genres").
		Preload("People.Person")

	if filter.SearchKeyword != "" {
		q = q.Where("movies.title ILIKE ?", "%"+filter.SearchKeyword+"%")
	}

	if len(filter.IDs) > 0 {
		q = q.Where("movies.id IN ?", filter.IDs)
	}

	if filter.Status != nil {
		q = q.Where("movies.status = ?", *filter.Status)
	}

	if filter.MinRating != nil {
		q = q.Where("movies.rating_avg >= ?", *filter.MinRating)
	}

	if filter.Country != nil && *filter.Country != "" {
		q = q.Where("movies.country = ?", *filter.Country)
	}

	if filter.GenreID != nil {
		q = q.Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
			Where("movie_genres.genre_id = ?", *filter.GenreID)
	}

	if filter.PersonID != nil {
		q = q.Joins("JOIN movie_people ON movie_people.movie_id = movies.id").
			Where("movie_people.person_id = ?", *filter.PersonID)

		if filter.PersonRole != nil {
			q = q.Where("movie_people.role = ?", *filter.PersonRole)
		}
		q = q.Group("movies.id")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	if opts != nil && opts.HasSorting() {
		direction := "ASC"
		if !opts.Sorting.Asc {
			direction = "DESC"
		}
		q = q.Order(fmt.Sprintf("movies.%s %s", opts.Sorting.Column, direction))
	} else {
		q = q.Order("movies.created_at DESC")
	}

	if opts != nil && opts.HasPagination() {
		q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	if err := q.Find(&movies).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}

	return movies, total, nil
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

func (r *movieRepository) UpdateRatingAggregate(ctx context.Context, movieID int64) error {
	query := `
        UPDATE movies
        SET 
            rating_avg = COALESCE((SELECT ROUND(AVG(score)::numeric, 2) FROM ratings WHERE movie_id = ?), 0.00),
            rating_count = (SELECT COUNT(*) FROM ratings WHERE movie_id = ?)
        WHERE id = ?
    `
	err := r.db.GetDB().WithContext(ctx).Exec(query, movieID, movieID, movieID).Error
	if err != nil {
		return fmt.Errorf("%w: %v", exception.ErrDatabase, err)
	}
	return nil
}

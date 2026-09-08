package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/enum"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/DB"
	"github.com/mehdi1443/cinema-booking-system/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	return &movieRepository{
		db: db,
	}
}

// -----------------------------------------------------------------------------
// Database Error Handling
// -----------------------------------------------------------------------------

func handleDBError(err error) error {
	if err == nil {
		return nil
	}

	// GORM translated errors.
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf(
			"duplicate database key: %w",
			exception.ErrDuplicateKey,
		)
	}

	// Check native PostgreSQL error.
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case "23505":
			// unique_violation
			return fmt.Errorf(
				"postgres unique violation: %w",
				exception.ErrDuplicateKey,
			)

		case "23503":
			// foreign_key_violation

			constraint := strings.ToLower(pgErr.ConstraintName)
			detail := strings.ToLower(pgErr.Detail)
			table := strings.ToLower(pgErr.TableName)

			// Genre relation.
			if strings.Contains(constraint, "genre") ||
				strings.Contains(detail, "movie_genres") ||
				strings.Contains(table, "movie_genres") {

				return fmt.Errorf(
					"genre foreign key violation: %w",
					exception.ErrGenreNotFound,
				)
			}

			// Person relation.
			if strings.Contains(constraint, "person") ||
				strings.Contains(detail, "movie_people") ||
				strings.Contains(table, "movie_people") {

				return fmt.Errorf(
					"person foreign key violation: %w",
					exception.ErrPersonNotFound,
				)
			}

			return fmt.Errorf(
				"unknown foreign key violation: %w",
				exception.ErrDatabase,
			)

		case "23502":
			// not_null_violation
			return fmt.Errorf(
				"database not null violation: %w",
				exception.ErrDatabase,
			)

		case "23514":
			// check_violation
			return fmt.Errorf(
				"database check violation: %w",
				exception.ErrDatabase,
			)

		case "40001":
			// serialization_failure
			return fmt.Errorf(
				"database serialization failure: %w",
				exception.ErrDatabase,
			)

		case "40P01":
			// deadlock_detected
			return fmt.Errorf(
				"database deadlock detected: %w",
				exception.ErrDatabase,
			)
		}
	}

	// Unknown database error.
	return fmt.Errorf(
		"database operation failed: %w",
		exception.ErrDatabase,
	)
}

// -----------------------------------------------------------------------------
// Create
// -----------------------------------------------------------------------------

func (r *movieRepository) Create(
	ctx context.Context,
	movie *entity.Movie,
	genreIDs []int64,
	peopleIDs map[int64]enum.MoviePersonRole,
) (*entity.Movie, error) {

	err := r.db.WithTransaction(func(txWrapper database.Database) error {

		tx := txWrapper.GetDB().WithContext(ctx)

		// Create movie without automatic relations.
		if err := tx.
			Omit("ID", "Genres", "People").
			Create(movie).
			Error; err != nil {

			return handleDBError(err)
		}

		// Genres.
		if len(genreIDs) > 0 {

			movieGenres := make(
				[]MovieGenre,
				0,
				len(genreIDs),
			)

			for _, genreID := range genreIDs {
				movieGenres = append(
					movieGenres,
					MovieGenre{
						MovieID: movie.ID,
						GenreID: genreID,
					},
				)
			}

			if err := tx.
				Create(&movieGenres).
				Error; err != nil {

				return handleDBError(err)
			}
		}

		// People.
		if len(peopleIDs) > 0 {

			moviePeople := make(
				[]entity.MoviePerson,
				0,
				len(peopleIDs),
			)

			for personID, role := range peopleIDs {

				moviePeople = append(
					moviePeople,
					entity.MoviePerson{
						MovieID:  movie.ID,
						PersonID: personID,
						Role:     role,
					},
				)
			}

			if err := tx.
				Create(&moviePeople).
				Error; err != nil {

				return handleDBError(err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, movie.ID)
}

// -----------------------------------------------------------------------------
// Get By ID
// -----------------------------------------------------------------------------

func (r *movieRepository) GetByID(
	ctx context.Context,
	id int64,
) (*entity.Movie, error) {

	var movie entity.Movie

	err := r.db.
		GetDB().
		WithContext(ctx).
		Preload("Genres").
		Preload("People.Person").
		First(&movie, id).
		Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.ErrMovieNotFound
		}

		return nil, handleDBError(err)
	}

	return &movie, nil
}

// -----------------------------------------------------------------------------
// Get Movies
// -----------------------------------------------------------------------------

func (r *movieRepository) GetMovies(
	ctx context.Context,
	filter DB.MovieFilter,
	opts *DB.QueryOptions,
) ([]entity.Movie, int64, error) {

	var movies []entity.Movie

	db := r.db.GetDB().WithContext(ctx)

	q := db.Model(&entity.Movie{})

	if filter.SearchKeyword != "" {
		q = q.Where(
			"movies.title ILIKE ?",
			"%"+filter.SearchKeyword+"%",
		)
	}

	if len(filter.IDs) > 0 {
		q = q.Where(
			"movies.id IN ?",
			filter.IDs,
		)
	}

	if filter.Status != nil {
		q = q.Where(
			"movies.status = ?",
			*filter.Status,
		)
	}

	if filter.MinRating != nil {
		q = q.Where(
			"movies.rating_avg >= ?",
			*filter.MinRating,
		)
	}

	if filter.MaxRating != nil {
		q = q.Where(
			"movies.rating_avg <= ?",
			*filter.MaxRating,
		)
	}

	if filter.Country != nil && *filter.Country != "" {
		q = q.Where(
			"movies.country = ?",
			*filter.Country,
		)
	}

	if filter.FromYear != nil {
		q = q.Where(
			"EXTRACT(YEAR FROM movies.release_date) >= ?",
			*filter.FromYear,
		)
	}

	if filter.ToYear != nil {
		q = q.Where(
			"EXTRACT(YEAR FROM movies.release_date) <= ?",
			*filter.ToYear,
		)
	}

	if filter.MinDuration != nil {
		q = q.Where(
			"movies.duration_minutes >= ?",
			*filter.MinDuration,
		)
	}

	if filter.MaxDuration != nil {
		q = q.Where(
			"movies.duration_minutes <= ?",
			*filter.MaxDuration,
		)
	}

	if filter.HasTrailer != nil {

		if *filter.HasTrailer {
			q = q.Where(
				"movies.trailer_url IS NOT NULL AND movies.trailer_url != ''",
			)
		} else {
			q = q.Where(
				"movies.trailer_url IS NULL OR movies.trailer_url = ''",
			)
		}
	}

	// Genres:
	// Movie must contain ALL requested genres.
	if len(filter.GenreIDs) > 0 {

		subQuery := db.
			Table("movie_genres").
			Select("movie_id").
			Where(
				"genre_id IN ?",
				filter.GenreIDs,
			).
			Group("movie_id").
			Having(
				"COUNT(DISTINCT genre_id) = ?",
				len(filter.GenreIDs),
			)

		q = q.Where(
			"movies.id IN (?)",
			subQuery,
		)
	}

	// People:
	// Movie must contain ALL requested people.
	if len(filter.PersonIDs) > 0 {

		subQuery := db.
			Table("movie_people").
			Select("movie_id").
			Where(
				"person_id IN ?",
				filter.PersonIDs,
			)

		if filter.PersonRole != nil {
			subQuery = subQuery.Where(
				"role = ?",
				*filter.PersonRole,
			)
		}

		subQuery = subQuery.
			Group("movie_id").
			Having(
				"COUNT(DISTINCT person_id) = ?",
				len(filter.PersonIDs),
			)

		q = q.Where(
			"movies.id IN (?)",
			subQuery,
		)
	}

	// Count.
	var total int64

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, handleDBError(err)
	}

	// Relations.
	q = q.
		Preload("Genres").
		Preload("People.Person")

	// Sorting.
	if opts != nil && opts.HasSorting() {

		direction := "ASC"

		if !opts.Sorting.Asc {
			direction = "DESC"
		}

		q = q.Order(
			fmt.Sprintf(
				"movies.%s %s",
				opts.Sorting.Column,
				direction,
			),
		)
	} else {
		q = q.Order("movies.created_at DESC")
	}

	// Pagination.
	if opts != nil && opts.HasPagination() {

		q = q.
			Offset(opts.Pagination.Offset).
			Limit(opts.Pagination.Limit)
	}

	// Execute.
	if err := q.Find(&movies).Error; err != nil {
		return nil, 0, handleDBError(err)
	}

	return movies, total, nil
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func (r *movieRepository) Update(
	ctx context.Context,
	movie *entity.Movie,
	genreIDs []int64,
	peopleIDs map[int64]enum.MoviePersonRole,
) (*entity.Movie, error) {

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

		result := tx.
			Model(&entity.Movie{}).
			Where("id = ?", movie.ID).
			Updates(updateData)

		if result.Error != nil {
			return handleDBError(result.Error)
		}

		if result.RowsAffected == 0 {
			return exception.ErrMovieNotFound
		}

		// ---------------------------------------------------------------------
		// Genres
		// nil  => don't touch genres
		// empty => remove all genres
		// ---------------------------------------------------------------------

		if genreIDs != nil {

			if err := tx.
				Where("movie_id = ?", movie.ID).
				Delete(&MovieGenre{}).
				Error; err != nil {

				return handleDBError(err)
			}

			if len(genreIDs) > 0 {

				movieGenres := make(
					[]MovieGenre,
					0,
					len(genreIDs),
				)

				for _, genreID := range genreIDs {
					movieGenres = append(
						movieGenres,
						MovieGenre{
							MovieID: movie.ID,
							GenreID: genreID,
						},
					)
				}

				if err := tx.
					Create(&movieGenres).
					Error; err != nil {

					return handleDBError(err)
				}
			}
		}

		// ---------------------------------------------------------------------
		// People
		// nil  => don't touch people
		// empty => remove all people
		// ---------------------------------------------------------------------

		if peopleIDs != nil {

			if err := tx.
				Where("movie_id = ?", movie.ID).
				Delete(&entity.MoviePerson{}).
				Error; err != nil {

				return handleDBError(err)
			}

			if len(peopleIDs) > 0 {

				moviePeople := make(
					[]entity.MoviePerson,
					0,
					len(peopleIDs),
				)

				for personID, role := range peopleIDs {
					moviePeople = append(
						moviePeople,
						entity.MoviePerson{
							MovieID:  movie.ID,
							PersonID: personID,
							Role:     role,
						},
					)
				}

				if err := tx.
					Create(&moviePeople).
					Error; err != nil {

					return handleDBError(err)
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

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func (r *movieRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	return r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var movie entity.Movie
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&movie, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrMovieNotFound
			}
			return handleDBError(err)
		}

		var showtimeCount int64
		if err := tx.Model(&entity.Showtime{}).
			Where("movie_id = ? AND starts_at > ? AND status IN ?", id, time.Now(), []enum.ShowtimeStatus{enum.ShowtimeStatusDraft, enum.ShowtimeStatusPublished}).
			Count(&showtimeCount).Error; err != nil {
			return handleDBError(err)
		}
		if showtimeCount > 0 {
			return exception.ErrResourceInUse
		}

		result := tx.Delete(&movie)
		if result.Error != nil {
			return handleDBError(result.Error)
		}
		if result.RowsAffected == 0 {
			return exception.ErrMovieNotFound
		}
		return nil
	})
}

// -----------------------------------------------------------------------------
// Rating Aggregate
// -----------------------------------------------------------------------------

func (r *movieRepository) UpdateRatingAggregate(
	ctx context.Context,
	movieID int64,
) error {

	query := `
		UPDATE movies
		SET
			rating_avg = COALESCE(
				(
					SELECT ROUND(
						AVG(score)::numeric,
						2
					)
					FROM ratings
					WHERE movie_id = ?
				),
				0.00
			),
			rating_count = (
				SELECT COUNT(*)
				FROM ratings
				WHERE movie_id = ?
			)
		WHERE id = ?
	`

	err := r.db.
		GetDB().
		WithContext(ctx).
		Exec(
			query,
			movieID,
			movieID,
			movieID,
		).
		Error

	if err != nil {
		return handleDBError(err)
	}

	return nil
}

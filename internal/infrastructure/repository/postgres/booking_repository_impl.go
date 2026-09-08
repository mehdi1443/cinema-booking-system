package postgres_repo

import (
	"context"
	"errors"
	"fmt"
	"math"
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

type bookingRepository struct {
	db database.Database
}

func NewBookingRepository(db database.Database) DB.BookingRepositoryInterface {
	return &bookingRepository{db: db}
}

func handleBookingDBError(err error) error {
	if err == nil {
		return nil
	}
	var appErr *exception.AppError
	if errors.As(err, &appErr) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrRecordNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("duplicate booking key: %w", exception.ErrDuplicateKey)
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		constraint := strings.ToLower(pgErr.ConstraintName)
		switch pgErr.Code {
		case "23P01":
			return fmt.Errorf("showtime exclusion constraint violation: %w", exception.ErrShowtimeOverlap)
		case "23505":
			switch {
			case strings.Contains(constraint, "uq_halls"):
				return fmt.Errorf("duplicate hall: %w", exception.ErrHallAlreadyExists)
			case strings.Contains(constraint, "uq_seats_hall"):
				return fmt.Errorf("duplicate seat position: %w", exception.ErrSeatAlreadyExists)
			default:
				return fmt.Errorf("duplicate booking record: %w", exception.ErrDuplicateKey)
			}
		case "23503":
			return mapBookingForeignKeyError(pgErr)
		case "23514":
			switch {
			case strings.Contains(constraint, "price"), strings.Contains(constraint, "total"):
				return fmt.Errorf("invalid booking amount: %w", exception.ErrInvalidAmount)
			case strings.Contains(constraint, "number"), strings.Contains(constraint, "time_range"):
				return fmt.Errorf("invalid booking data: %w", exception.ErrInvalidInput)
			default:
				return fmt.Errorf("booking check constraint violation: %w", exception.ErrDatabase)
			}
		case "23502", "40001", "40P01":
			return fmt.Errorf("booking constraint or transaction failure: %w", exception.ErrDatabase)
		}
	}
	return fmt.Errorf("booking database operation failed: %w", exception.ErrDatabase)
}

func mapBookingForeignKeyError(pgErr *pgconn.PgError) error {
	constraint := strings.ToLower(pgErr.ConstraintName)
	detail := strings.ToLower(pgErr.Detail)
	reference := constraint + " " + detail

	switch {
	case strings.Contains(reference, "halls_cinema_id"), strings.Contains(detail, "(cinema_id)"):
		return fmt.Errorf("cinema foreign key violation: %w", exception.ErrCinemaNotFound)
	case strings.Contains(reference, "showtimes_movie_id"), strings.Contains(detail, "(movie_id)"):
		return fmt.Errorf("movie foreign key violation: %w", exception.ErrMovieNotFound)
	case strings.Contains(reference, "reservations_user_id"), strings.Contains(detail, "(user_id)"):
		return fmt.Errorf("user foreign key violation: %w", exception.ErrUserNotFound)
	case strings.Contains(reference, "seats_hall_id"), strings.Contains(reference, "showtimes_hall_id"), strings.Contains(detail, "(hall_id)"):
		return fmt.Errorf("hall foreign key violation: %w", exception.ErrHallNotFound)
	case strings.Contains(reference, "reservations_showtime_id"), strings.Contains(reference, "showtime_seats_showtime_id"), strings.Contains(detail, "(showtime_id)"):
		return fmt.Errorf("showtime foreign key violation: %w", exception.ErrShowtimeNotFound)
	case strings.Contains(reference, "showtime_seats_reservation_id"), strings.Contains(reference, "reservation_items_reservation_id"), strings.Contains(detail, "(reservation_id)"):
		return fmt.Errorf("reservation foreign key violation: %w", exception.ErrReservationNotFound)
	case strings.Contains(reference, "showtime_seats_seat_id"), strings.Contains(reference, "reservation_items_showtime_seat_id"), strings.Contains(detail, "(seat_id)"), strings.Contains(detail, "(showtime_seat_id)"):
		return fmt.Errorf("seat foreign key violation: %w", exception.ErrSeatNotFound)
	default:
		return fmt.Errorf("unknown booking foreign key violation: %w", exception.ErrDatabase)
	}
}

func (r *bookingRepository) CreateHall(ctx context.Context, hall *entity.Hall) (*entity.Hall, error) {
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var cinema entity.Cinema
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&cinema, hall.CinemaID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrCinemaNotFound
			}
			return handleBookingDBError(err)
		}
		return handleBookingDBError(tx.Create(hall).Error)
	})
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return r.GetHall(ctx, hall.ID)
}

func (r *bookingRepository) GetHall(ctx context.Context, id int64) (*entity.Hall, error) {
	var hall entity.Hall
	err := r.db.GetDB().WithContext(ctx).
		Preload("Cinema").Preload("Seats", func(db *gorm.DB) *gorm.DB {
		return db.Order("row_label ASC, number ASC")
	}).First(&hall, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrHallNotFound
	}
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return &hall, nil
}

func (r *bookingRepository) ListHalls(ctx context.Context, cinemaID int64) ([]entity.Hall, error) {
	var halls []entity.Hall
	err := r.db.GetDB().WithContext(ctx).Where("cinema_id = ?", cinemaID).
		Preload("Seats", func(db *gorm.DB) *gorm.DB { return db.Order("row_label ASC, number ASC") }).
		Order("name ASC").Find(&halls).Error
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return halls, nil
}

func (r *bookingRepository) UpdateHall(ctx context.Context, hall *entity.Hall) (*entity.Hall, error) {
	result := r.db.GetDB().WithContext(ctx).Model(&entity.Hall{}).Where("id = ?", hall.ID).
		Updates(map[string]any{"name": hall.Name, "status": hall.Status, "updated_at": time.Now()})
	if result.Error != nil {
		return nil, handleBookingDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exception.ErrHallNotFound
	}
	return r.GetHall(ctx, hall.ID)
}

func (r *bookingRepository) DeleteHall(ctx context.Context, id int64) error {
	return r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var hall entity.Hall
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&hall, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrHallNotFound
			}
			return handleBookingDBError(err)
		}
		var count int64
		if err := tx.Model(&entity.Showtime{}).Where("hall_id = ? AND starts_at > ? AND status IN ?", id, time.Now(), []enum.ShowtimeStatus{enum.ShowtimeStatusDraft, enum.ShowtimeStatusPublished}).Count(&count).Error; err != nil {
			return handleBookingDBError(err)
		}
		if count > 0 {
			return exception.ErrResourceInUse
		}
		result := tx.Delete(&entity.Hall{}, id)
		if result.Error != nil {
			return handleBookingDBError(result.Error)
		}
		if result.RowsAffected == 0 {
			return exception.ErrHallNotFound
		}
		return nil
	})
}

func (r *bookingRepository) CreateSeats(ctx context.Context, hallID int64, seats []entity.Seat) ([]entity.Seat, error) {
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var hall entity.Hall
		if err := tx.First(&hall, hallID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrHallNotFound
			}
			return handleBookingDBError(err)
		}
		for i := range seats {
			seats[i].HallID = hallID
		}
		if err := tx.Create(&seats).Error; err != nil {
			return handleBookingDBError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *bookingRepository) GetSeat(ctx context.Context, id int64) (*entity.Seat, error) {
	var seat entity.Seat
	err := r.db.GetDB().WithContext(ctx).Preload("Hall.Cinema").First(&seat, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrSeatNotFound
	}
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return &seat, nil
}

func (r *bookingRepository) ListSeats(ctx context.Context, hallID int64) ([]entity.Seat, error) {
	var seats []entity.Seat
	if err := r.db.GetDB().WithContext(ctx).Where("hall_id = ?", hallID).Order("row_label ASC, number ASC").Find(&seats).Error; err != nil {
		return nil, handleBookingDBError(err)
	}
	return seats, nil
}

func (r *bookingRepository) UpdateSeat(ctx context.Context, seat *entity.Seat) (*entity.Seat, error) {
	result := r.db.GetDB().WithContext(ctx).Model(&entity.Seat{}).Where("id = ?", seat.ID).
		Updates(map[string]any{"row_label": seat.RowLabel, "number": seat.Number, "seat_type": seat.Type, "status": seat.Status, "updated_at": time.Now()})
	if result.Error != nil {
		return nil, handleBookingDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exception.ErrSeatNotFound
	}
	return r.GetSeat(ctx, seat.ID)
}

func (r *bookingRepository) DeleteSeat(ctx context.Context, id int64) error {
	result := r.db.GetDB().WithContext(ctx).Delete(&entity.Seat{}, id)
	if result.Error != nil {
		return handleBookingDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return exception.ErrSeatNotFound
	}
	return nil
}

func (r *bookingRepository) IsHallOwnedBy(ctx context.Context, hallID, userID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).Table("halls").
		Joins("JOIN cinemas ON cinemas.id = halls.cinema_id").
		Where("halls.id = ? AND halls.deleted_at IS NULL AND cinemas.owner_user_id = ? AND cinemas.deleted_at IS NULL", hallID, userID).Count(&count).Error
	return count > 0, handleBookingDBError(err)
}

func (r *bookingRepository) IsSeatOwnedBy(ctx context.Context, seatID, userID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).Table("seats").
		Joins("JOIN halls ON halls.id = seats.hall_id").Joins("JOIN cinemas ON cinemas.id = halls.cinema_id").
		Where("seats.id = ? AND seats.deleted_at IS NULL AND halls.deleted_at IS NULL AND cinemas.owner_user_id = ? AND cinemas.deleted_at IS NULL", seatID, userID).Count(&count).Error
	return count > 0, handleBookingDBError(err)
}

func (r *bookingRepository) CreateShowtime(ctx context.Context, showtime *entity.Showtime, prices map[enum.SeatType]int64) (*entity.Showtime, error) {
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var hall entity.Hall
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&hall, showtime.HallID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrHallNotFound
			}
			return handleBookingDBError(err)
		}
		if hall.Status != enum.HallStatusActive {
			return exception.ErrShowtimeNotBookable
		}
		var movie entity.Movie
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&movie, showtime.MovieID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrMovieNotFound
			}
			return handleBookingDBError(err)
		}
		if movie.DurationMinutes <= 0 || movie.Status == enum.MovieInactive || movie.Status == enum.MovieEnded {
			return exception.ErrShowtimeNotBookable
		}
		showtime.EndsAt = showtime.StartsAt.Add(time.Duration(movie.DurationMinutes) * time.Minute)
		if err := tx.Create(showtime).Error; err != nil {
			return handleBookingDBError(err)
		}
		var seats []entity.Seat
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Where("hall_id = ? AND status = ?", showtime.HallID, enum.SeatStatusActive).Order("id ASC").Find(&seats).Error; err != nil {
			return handleBookingDBError(err)
		}
		if len(seats) == 0 {
			return exception.ErrNoActiveSeats
		}
		snapshots := make([]entity.ShowtimeSeat, 0, len(seats))
		for i := range seats {
			price, ok := prices[seats[i].Type]
			if !ok || price <= 0 {
				return exception.ErrInvalidAmount
			}
			snapshots = append(snapshots, entity.ShowtimeSeat{ShowtimeID: showtime.ID, SeatID: seats[i].ID, RowLabel: seats[i].RowLabel, SeatNumber: seats[i].Number, SeatType: seats[i].Type, Price: price, Status: enum.ShowtimeSeatAvailable})
		}
		return handleBookingDBError(tx.Create(&snapshots).Error)
	})
	if err != nil {
		return nil, err
	}
	return r.GetShowtime(ctx, showtime.ID)
}

func (r *bookingRepository) GetShowtime(ctx context.Context, id int64) (*entity.Showtime, error) {
	var showtime entity.Showtime
	err := r.db.GetDB().WithContext(ctx).Preload("Hall.Cinema").Preload("Movie").Preload("Seats").First(&showtime, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrShowtimeNotFound
	}
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return &showtime, nil
}

func (r *bookingRepository) ListShowtimes(ctx context.Context, filter DB.ShowtimeFilter, opts *DB.QueryOptions) ([]entity.Showtime, int64, error) {
	var showtimes []entity.Showtime
	query := r.db.GetDB().WithContext(ctx).Model(&entity.Showtime{})
	hallsJoined := false
	if filter.PublicOnly {
		query = query.
			Joins("JOIN halls ON halls.id = showtimes.hall_id AND halls.deleted_at IS NULL").
			Joins("JOIN cinemas ON cinemas.id = halls.cinema_id AND cinemas.deleted_at IS NULL").
			Joins("JOIN movies ON movies.id = showtimes.movie_id AND movies.deleted_at IS NULL").
			Where("halls.status = ? AND cinemas.status = ?", enum.HallStatusActive, enum.CinemaStatusActive).
			Where("movies.status NOT IN ?", []enum.MovieStatus{enum.MovieInactive, enum.MovieEnded})
		hallsJoined = true
	}
	if filter.CinemaID != nil {
		if !hallsJoined {
			query = query.Joins("JOIN halls ON halls.id = showtimes.hall_id AND halls.deleted_at IS NULL")
		}
		query = query.Where("halls.cinema_id = ?", *filter.CinemaID)
	}
	if filter.HallID != nil {
		query = query.Where("showtimes.hall_id = ?", *filter.HallID)
	}
	if filter.MovieID != nil {
		query = query.Where("showtimes.movie_id = ?", *filter.MovieID)
	}
	if filter.Status != nil {
		query = query.Where("showtimes.status = ?", *filter.Status)
	}
	if filter.From != nil {
		query = query.Where("showtimes.starts_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("showtimes.starts_at < ?", *filter.To)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, handleBookingDBError(err)
	}
	query = query.Preload("Hall.Cinema").Preload("Movie").Preload("Seats").Order("showtimes.starts_at ASC")
	if opts != nil && opts.HasPagination() {
		query = query.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}
	if err := query.Find(&showtimes).Error; err != nil {
		return nil, 0, handleBookingDBError(err)
	}
	return showtimes, total, nil
}

func (r *bookingRepository) UpdateShowtime(ctx context.Context, showtime *entity.Showtime) (*entity.Showtime, error) {
	result := r.db.GetDB().WithContext(ctx).Model(&entity.Showtime{}).Where("id = ?", showtime.ID).
		Updates(map[string]any{"starts_at": showtime.StartsAt, "ends_at": showtime.EndsAt, "status": showtime.Status, "updated_at": time.Now()})
	if result.Error != nil {
		return nil, handleBookingDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exception.ErrShowtimeNotFound
	}
	return r.GetShowtime(ctx, showtime.ID)
}

func (r *bookingRepository) DeleteShowtime(ctx context.Context, id int64) error {
	return r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&showtime, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handleBookingDBError(err)
		}
		var count int64
		if err := tx.Model(&entity.ShowtimeSeat{}).Where("showtime_id = ? AND status IN ?", id, []enum.ShowtimeSeatStatus{enum.ShowtimeSeatHeld, enum.ShowtimeSeatSold}).Count(&count).Error; err != nil {
			return handleBookingDBError(err)
		}
		if count > 0 {
			return exception.ErrResourceInUse
		}
		result := tx.Delete(&entity.Showtime{}, id)
		if result.Error != nil {
			return handleBookingDBError(result.Error)
		}
		if result.RowsAffected == 0 {
			return exception.ErrShowtimeNotFound
		}
		return nil
	})
}

func (r *bookingRepository) UpdateShowtimePrices(ctx context.Context, showtimeID int64, prices map[int64]int64) error {
	ids := make([]int64, 0, len(prices))
	for id := range prices {
		ids = append(ids, id)
	}
	return r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&showtime, showtimeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handleBookingDBError(err)
		}
		if !showtime.StartsAt.After(time.Now()) || showtime.Status == enum.ShowtimeStatusCancelled || showtime.Status == enum.ShowtimeStatusCompleted {
			return exception.ErrShowtimeNotBookable
		}
		var seats []entity.ShowtimeSeat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("showtime_id = ? AND id IN ?", showtimeID, ids).Order("id ASC").Find(&seats).Error; err != nil {
			return handleBookingDBError(err)
		}
		if len(seats) != len(prices) {
			return exception.ErrInvalidSeatSelection
		}
		for i := range seats {
			if seats[i].Status == enum.ShowtimeSeatHeld || seats[i].Status == enum.ShowtimeSeatSold {
				return exception.ErrResourceInUse
			}
			if err := tx.Model(&entity.ShowtimeSeat{}).Where("id = ?", seats[i].ID).Updates(map[string]any{"price": prices[seats[i].ID], "updated_at": time.Now()}).Error; err != nil {
				return handleBookingDBError(err)
			}
		}
		return nil
	})
}

func (r *bookingRepository) ListShowtimeSeats(ctx context.Context, showtimeID int64) ([]entity.ShowtimeSeat, error) {
	var count int64
	if err := r.db.GetDB().WithContext(ctx).Model(&entity.Showtime{}).Where("id = ?", showtimeID).Count(&count).Error; err != nil {
		return nil, handleBookingDBError(err)
	}
	if count == 0 {
		return nil, exception.ErrShowtimeNotFound
	}
	var seats []entity.ShowtimeSeat
	if err := r.db.GetDB().WithContext(ctx).Where("showtime_id = ?", showtimeID).Order("row_label ASC, seat_number ASC").Find(&seats).Error; err != nil {
		return nil, handleBookingDBError(err)
	}
	return seats, nil
}

func (r *bookingRepository) IsShowtimeOwnedBy(ctx context.Context, showtimeID, userID int64) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).Table("showtimes").Joins("JOIN halls ON halls.id = showtimes.hall_id").Joins("JOIN cinemas ON cinemas.id = halls.cinema_id").
		Where("showtimes.id = ? AND showtimes.deleted_at IS NULL AND halls.deleted_at IS NULL AND cinemas.owner_user_id = ? AND cinemas.deleted_at IS NULL", showtimeID, userID).Count(&count).Error
	return count > 0, handleBookingDBError(err)
}

func (r *bookingRepository) CreateReservation(ctx context.Context, userID, showtimeID int64, showtimeSeatIDs []int64, expiresAt time.Time) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		now := time.Now()
		if err := expireReservations(tx, now); err != nil {
			return handleBookingDBError(err)
		}
		var showtime entity.Showtime
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&showtime, showtimeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrShowtimeNotFound
			}
			return handleBookingDBError(err)
		}
		if showtime.Status != enum.ShowtimeStatusPublished || !showtime.StartsAt.After(now) {
			return exception.ErrShowtimeNotBookable
		}
		if !expiresAt.Before(showtime.StartsAt) {
			expiresAt = showtime.StartsAt
		}
		var seats []entity.ShowtimeSeat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("showtime_id = ? AND id IN ?", showtimeID, showtimeSeatIDs).Order("id ASC").Find(&seats).Error; err != nil {
			return handleBookingDBError(err)
		}
		if len(seats) != len(showtimeSeatIDs) {
			return exception.ErrInvalidSeatSelection
		}
		var total int64
		for i := range seats {
			if seats[i].Status != enum.ShowtimeSeatAvailable {
				return exception.ErrSeatUnavailable
			}
			if seats[i].Price <= 0 || total > math.MaxInt64-seats[i].Price {
				return exception.ErrInvalidAmount
			}
			total += seats[i].Price
		}
		reservation = entity.Reservation{UserID: userID, ShowtimeID: showtimeID, Status: enum.ReservationPendingPayment, ExpiresAt: expiresAt, TotalAmount: total}
		if err := tx.Create(&reservation).Error; err != nil {
			return handleBookingDBError(err)
		}
		items := make([]entity.ReservationItem, 0, len(seats))
		for i := range seats {
			items = append(items, entity.ReservationItem{ReservationID: reservation.ID, ShowtimeSeatID: seats[i].ID, UnitPrice: seats[i].Price})
			if err := tx.Model(&entity.ShowtimeSeat{}).Where("id = ?", seats[i].ID).Updates(map[string]any{"status": enum.ShowtimeSeatHeld, "reservation_id": reservation.ID, "held_until": expiresAt, "updated_at": now}).Error; err != nil {
				return handleBookingDBError(err)
			}
		}
		if err := tx.Create(&items).Error; err != nil {
			return handleBookingDBError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetReservation(ctx, reservation.ID)
}

func (r *bookingRepository) GetReservation(ctx context.Context, id int64) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.db.GetDB().WithContext(ctx).Preload("Showtime.Hall.Cinema").Preload("Showtime.Movie").Preload("Items.ShowtimeSeat").First(&reservation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrReservationNotFound
	}
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return &reservation, nil
}

func (r *bookingRepository) ListUserReservations(ctx context.Context, userID int64) ([]entity.Reservation, error) {
	var reservations []entity.Reservation
	err := r.db.GetDB().WithContext(ctx).Where("user_id = ?", userID).Preload("Showtime.Hall.Cinema").Preload("Showtime.Movie").Preload("Items.ShowtimeSeat").Order("created_at DESC").Find(&reservations).Error
	if err != nil {
		return nil, handleBookingDBError(err)
	}
	return reservations, nil
}

func (r *bookingRepository) CancelReservation(ctx context.Context, id int64) (*entity.Reservation, error) {
	err := r.db.WithTransaction(func(txDB database.Database) error {
		tx := txDB.GetDB().WithContext(ctx)
		var reservation entity.Reservation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&reservation, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.ErrReservationNotFound
			}
			return handleBookingDBError(err)
		}
		if reservation.Status == enum.ReservationCancelled || reservation.Status == enum.ReservationExpired {
			return nil
		}
		if reservation.Status != enum.ReservationPendingPayment {
			return exception.ErrResourceInUse
		}
		now := time.Now()
		status := enum.ReservationCancelled
		if !reservation.ExpiresAt.After(now) {
			status = enum.ReservationExpired
		}
		if err := tx.Model(&entity.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": status, "updated_at": now}).Error; err != nil {
			return handleBookingDBError(err)
		}
		if err := tx.Model(&entity.ShowtimeSeat{}).Where("reservation_id = ? AND status = ?", id, enum.ShowtimeSeatHeld).Updates(map[string]any{"status": enum.ShowtimeSeatAvailable, "reservation_id": nil, "held_until": nil, "updated_at": now}).Error; err != nil {
			return handleBookingDBError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetReservation(ctx, id)
}

func (r *bookingRepository) ExpireReservations(ctx context.Context, now time.Time) error {
	return handleBookingDBError(expireReservations(r.db.GetDB().WithContext(ctx), now))
}

func expireReservations(db *gorm.DB, now time.Time) error {
	query := `WITH expired AS (
    UPDATE reservations
       SET status = 'EXPIRED', updated_at = ?
     WHERE status = 'PENDING_PAYMENT' AND expires_at <= ?
     RETURNING id
)
UPDATE showtime_seats
   SET status = 'AVAILABLE', reservation_id = NULL, held_until = NULL, updated_at = ?
 WHERE status = 'HELD' AND held_until <= ?`
	return db.Exec(query, now, now, now, now).Error
}

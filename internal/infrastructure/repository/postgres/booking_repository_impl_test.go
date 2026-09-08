package postgres_repo

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/exception"
)

func TestHandleBookingDBErrorMapsForeignKeys(t *testing.T) {
	tests := []struct {
		constraint string
		expected   error
	}{
		{constraint: "halls_cinema_id_fkey", expected: exception.ErrCinemaNotFound},
		{constraint: "seats_hall_id_fkey", expected: exception.ErrHallNotFound},
		{constraint: "showtimes_movie_id_fkey", expected: exception.ErrMovieNotFound},
		{constraint: "reservations_user_id_fkey", expected: exception.ErrUserNotFound},
		{constraint: "reservations_showtime_id_fkey", expected: exception.ErrShowtimeNotFound},
		{constraint: "showtime_seats_reservation_id_fkey", expected: exception.ErrReservationNotFound},
		{constraint: "showtime_seats_seat_id_fkey", expected: exception.ErrSeatNotFound},
	}

	for _, test := range tests {
		t.Run(test.constraint, func(t *testing.T) {
			err := handleBookingDBError(&pgconn.PgError{Code: "23503", ConstraintName: test.constraint})
			if !errors.Is(err, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func TestHandleBookingDBErrorMapsBookingConstraints(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		constraint string
		expected   error
	}{
		{name: "overlap", code: "23P01", constraint: "ex_showtimes_hall_no_overlap", expected: exception.ErrShowtimeOverlap},
		{name: "duplicate hall", code: "23505", constraint: "uq_halls_cinema_name_active", expected: exception.ErrHallAlreadyExists},
		{name: "duplicate seat", code: "23505", constraint: "uq_seats_hall_position_active", expected: exception.ErrSeatAlreadyExists},
		{name: "invalid price", code: "23514", constraint: "chk_showtime_seats_price_positive", expected: exception.ErrInvalidAmount},
		{name: "invalid time range", code: "23514", constraint: "chk_showtimes_time_range", expected: exception.ErrInvalidInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := handleBookingDBError(&pgconn.PgError{Code: test.code, ConstraintName: test.constraint})
			if !errors.Is(err, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

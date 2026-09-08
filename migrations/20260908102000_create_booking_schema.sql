-- +goose Up
CREATE TYPE hall_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE seat_type AS ENUM ('STANDARD', 'VIP', 'WHEELCHAIR');
CREATE TYPE seat_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE showtime_status AS ENUM ('DRAFT', 'PUBLISHED', 'CANCELLED', 'COMPLETED');
CREATE TYPE showtime_seat_status AS ENUM ('AVAILABLE', 'HELD', 'SOLD', 'BLOCKED');
CREATE TYPE reservation_status AS ENUM ('PENDING_PAYMENT', 'PAID', 'EXPIRED', 'CANCELLED');

CREATE TABLE halls (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cinema_id BIGINT NOT NULL REFERENCES cinemas(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    status hall_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX uq_halls_cinema_name_active
    ON halls (cinema_id, LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_halls_cinema_id ON halls(cinema_id);

CREATE TABLE seats (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hall_id BIGINT NOT NULL REFERENCES halls(id) ON DELETE RESTRICT,
    row_label VARCHAR(10) NOT NULL,
    number INT NOT NULL,
    seat_type seat_type NOT NULL DEFAULT 'STANDARD',
    status seat_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_seats_number_positive CHECK (number > 0)
);

CREATE UNIQUE INDEX uq_seats_hall_position_active
    ON seats (hall_id, UPPER(row_label), number) WHERE deleted_at IS NULL;
CREATE INDEX idx_seats_hall_id ON seats(hall_id);

CREATE TABLE showtimes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hall_id BIGINT NOT NULL REFERENCES halls(id) ON DELETE RESTRICT,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE RESTRICT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    status showtime_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_showtimes_time_range CHECK (ends_at > starts_at)
);

CREATE INDEX idx_showtimes_hall_start ON showtimes(hall_id, starts_at);
CREATE INDEX idx_showtimes_movie_start ON showtimes(movie_id, starts_at);

CREATE EXTENSION IF NOT EXISTS btree_gist;
ALTER TABLE showtimes ADD CONSTRAINT ex_showtimes_hall_no_overlap
    EXCLUDE USING gist (
        hall_id WITH =,
        tstzrange(starts_at, ends_at, '[)') WITH &&
    ) WHERE (deleted_at IS NULL AND status IN ('DRAFT'::showtime_status, 'PUBLISHED'::showtime_status));

CREATE TABLE reservations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    showtime_id BIGINT NOT NULL REFERENCES showtimes(id) ON DELETE RESTRICT,
    status reservation_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    expires_at TIMESTAMPTZ NOT NULL,
    total_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_reservations_total_positive CHECK (total_amount > 0)
);

CREATE INDEX idx_reservations_user_created ON reservations(user_id, created_at DESC);
CREATE INDEX idx_reservations_pending_expiry ON reservations(expires_at)
    WHERE status = 'PENDING_PAYMENT';

CREATE TABLE showtime_seats (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    showtime_id BIGINT NOT NULL REFERENCES showtimes(id) ON DELETE RESTRICT,
    seat_id BIGINT NOT NULL REFERENCES seats(id) ON DELETE RESTRICT,
    row_label VARCHAR(10) NOT NULL,
    seat_number INT NOT NULL,
    seat_type seat_type NOT NULL,
    price BIGINT NOT NULL,
    status showtime_seat_status NOT NULL DEFAULT 'AVAILABLE',
    reservation_id BIGINT REFERENCES reservations(id) ON DELETE SET NULL,
    held_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_showtime_seats_showtime_seat UNIQUE (showtime_id, seat_id),
    CONSTRAINT chk_showtime_seats_price_positive CHECK (price > 0),
    CONSTRAINT chk_showtime_seats_hold_fields CHECK (
        (status = 'HELD' AND reservation_id IS NOT NULL AND held_until IS NOT NULL)
        OR status <> 'HELD'
    )
);

CREATE INDEX idx_showtime_seats_availability ON showtime_seats(showtime_id, status);
CREATE INDEX idx_showtime_seats_reservation ON showtime_seats(reservation_id);
CREATE INDEX idx_showtime_seats_held_until ON showtime_seats(held_until) WHERE status = 'HELD';

CREATE TABLE reservation_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reservation_id BIGINT NOT NULL REFERENCES reservations(id) ON DELETE RESTRICT,
    showtime_seat_id BIGINT NOT NULL REFERENCES showtime_seats(id) ON DELETE RESTRICT,
    unit_price BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_reservation_items_seat UNIQUE (reservation_id, showtime_seat_id),
    CONSTRAINT chk_reservation_items_price_positive CHECK (unit_price > 0)
);

CREATE INDEX idx_reservation_items_reservation ON reservation_items(reservation_id);

-- +goose Down
DROP TABLE IF EXISTS reservation_items;
DROP TABLE IF EXISTS showtime_seats;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS showtimes;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS halls;

DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS showtime_seat_status;
DROP TYPE IF EXISTS showtime_status;
DROP TYPE IF EXISTS seat_status;
DROP TYPE IF EXISTS seat_type;
DROP TYPE IF EXISTS hall_status;

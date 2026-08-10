-- +goose Up
-- =====================================================
-- ENUM TYPES
-- =====================================================
CREATE TYPE cinema_status AS ENUM (
    'ACTIVE',
    'INACTIVE'
);

-- =====================================================
-- LOCATION
-- =====================================================
CREATE TABLE provinces (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE cities (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    province_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,

    CONSTRAINT uq_cities_province_name
       UNIQUE (province_id, name),

    CONSTRAINT fk_cities_province
       FOREIGN KEY (province_id)
            REFERENCES provinces(id)
            ON DELETE RESTRICT
);

-- =====================================================
-- CINEMAS & FACILITIES
-- =====================================================
CREATE TABLE cinemas (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    city_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT NOT NULL,
    latitude DECIMAL(9, 6),
    longitude DECIMAL(9, 6),
    phone VARCHAR(20),
    image_url_low TEXT,
    image_url_high TEXT,
    status cinema_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_cinemas_city
        FOREIGN KEY (city_id)
            REFERENCES cities(id)
            ON DELETE RESTRICT
);

CREATE TABLE facilities (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE cinema_facilities (
    cinema_id BIGINT NOT NULL,
    facility_id BIGINT NOT NULL,

    PRIMARY KEY (cinema_id, facility_id),

    CONSTRAINT fk_cinema_facilities_cinema
        FOREIGN KEY (cinema_id)
            REFERENCES cinemas(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_cinema_facilities_facility
        FOREIGN KEY (facility_id)
            REFERENCES facilities(id)
            ON DELETE CASCADE
);

-- =====================================================
-- INDEXES
-- =====================================================
CREATE INDEX idx_cities_province_id ON cities(province_id);
CREATE INDEX idx_cinemas_city_id ON cinemas(city_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cinemas_city_id;
DROP INDEX IF EXISTS idx_cities_province_id;

DROP TABLE IF EXISTS cinema_facilities;
DROP TABLE IF EXISTS facilities;
DROP TABLE IF EXISTS cinemas;
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS provinces;

DROP TYPE IF EXISTS cinema_status;

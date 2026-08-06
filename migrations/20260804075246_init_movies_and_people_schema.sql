-- +goose Up
-- =====================================================
-- 1. ENUM TYPES
-- =====================================================

CREATE TYPE movie_status AS ENUM (
    'UPCOMING',
    'NOW_SHOWING',
    'ENDED',
    'INACTIVE'
);

CREATE TYPE movie_person_role AS ENUM (
    'ACTOR',
    'DIRECTOR',
    'WRITER'
);

-- =====================================================
-- 2. TABLES
-- =====================================================

CREATE TABLE movies (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT NOT NULL,
    release_date DATE,
    country VARCHAR(100),
    poster_url_low TEXT,
    poster_url_high TEXT,
    trailer_url TEXT,
    status movie_status NOT NULL DEFAULT 'UPCOMING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_movies_duration CHECK (duration_minutes > 0)
);



CREATE TABLE genres (
   id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
   name VARCHAR(100) NOT NULL UNIQUE,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



CREATE TABLE movie_genres (
   movie_id BIGINT NOT NULL,
   genre_id BIGINT NOT NULL,

   PRIMARY KEY (movie_id, genre_id),

   CONSTRAINT fk_movie_genres_movie
      FOREIGN KEY (movie_id)
      REFERENCES movies(id)
      ON DELETE CASCADE,

   CONSTRAINT fk_movie_genres_genre
      FOREIGN KEY (genre_id)
      REFERENCES genres(id)
      ON DELETE CASCADE
);



CREATE TABLE people (
   id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   biography TEXT,
   image_url_low TEXT,
   image_url_high TEXT,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



CREATE TABLE movie_people (
   movie_id BIGINT NOT NULL,
   person_id BIGINT NOT NULL,
   role movie_person_role NOT NULL,

   PRIMARY KEY (movie_id, person_id, role),

   CONSTRAINT fk_movie_people_movie
       FOREIGN KEY (movie_id)
          REFERENCES movies(id)
          ON DELETE CASCADE,

   CONSTRAINT fk_movie_people_person
       FOREIGN KEY (person_id)
          REFERENCES people(id)
          ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS movie_people;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS movie_genres;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS movies;

DROP TYPE IF EXISTS movie_person_role;
DROP TYPE IF EXISTS movie_status;
-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS ratings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    movie_id BIGINT,
    cinema_id BIGINT,
    score INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_ratings_target
    CHECK (
        (movie_id IS NOT NULL AND cinema_id IS NULL)
        OR
        (movie_id IS NULL AND cinema_id IS NOT NULL)
    ),

    CONSTRAINT chk_ratings_score
    CHECK (score BETWEEN 1 AND 5),

    CONSTRAINT uq_ratings_movie
    UNIQUE (user_id, movie_id),

    CONSTRAINT uq_ratings_cinema
    UNIQUE (user_id, cinema_id),

    CONSTRAINT fk_ratings_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_ratings_movie
    FOREIGN KEY (movie_id)
    REFERENCES movies(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_ratings_cinema
    FOREIGN KEY (cinema_id)
    REFERENCES cinemas(id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_ratings_movie_id ON ratings(movie_id);
CREATE INDEX IF NOT EXISTS idx_ratings_cinema_id ON ratings(cinema_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON ratings(user_id);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS ratings;
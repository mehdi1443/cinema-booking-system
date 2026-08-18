-- +goose Up
CREATE TYPE review_status AS ENUM (
    'PENDING',
    'PUBLISHED',
    'REJECTED',
    'SPOILER'
);

CREATE TYPE vote_type AS ENUM (
    'LIKE',
    'DISLIKE'
);

CREATE TABLE reviews (
                         id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                         user_id BIGINT NOT NULL,
                         movie_id BIGINT,
                         cinema_id BIGINT,
                         parent_id BIGINT,
                         body TEXT NOT NULL,
                         likes_count INT NOT NULL DEFAULT 0,
                         dislikes_count INT NOT NULL DEFAULT 0,
                         status review_status NOT NULL DEFAULT 'PENDING',
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                         CONSTRAINT chk_reviews_target
                             CHECK (
                                 (movie_id IS NOT NULL AND cinema_id IS NULL)
                                     OR
                                 (movie_id IS NULL AND cinema_id IS NOT NULL)
                                 ),

                         CONSTRAINT fk_reviews_parent
                             FOREIGN KEY (parent_id)
                                 REFERENCES reviews(id)
                                 ON DELETE CASCADE,

                         CONSTRAINT fk_reviews_user
                             FOREIGN KEY (user_id)
                                 REFERENCES users(id)
                                 ON DELETE CASCADE,

                         CONSTRAINT fk_reviews_movie
                             FOREIGN KEY (movie_id)
                                 REFERENCES movies(id)
                                 ON DELETE CASCADE,

                         CONSTRAINT fk_reviews_cinema
                             FOREIGN KEY (cinema_id)
                                 REFERENCES cinemas(id)
                                 ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_reviews_movie ON reviews (user_id, movie_id) WHERE parent_id IS NULL;
CREATE UNIQUE INDEX uq_reviews_cinema ON reviews (user_id, cinema_id) WHERE parent_id IS NULL;

CREATE INDEX idx_reviews_movie_id ON reviews(movie_id);
CREATE INDEX idx_reviews_cinema_id ON reviews(cinema_id);
CREATE INDEX idx_reviews_parent_id ON reviews(parent_id);

CREATE TABLE review_votes (
                              user_id BIGINT NOT NULL,
                              review_id BIGINT NOT NULL,
                              vote vote_type NOT NULL,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                              PRIMARY KEY (user_id, review_id),

                              CONSTRAINT fk_review_votes_user
                                  FOREIGN KEY (user_id)
                                      REFERENCES users(id)
                                      ON DELETE CASCADE,

                              CONSTRAINT fk_review_votes_review
                                  FOREIGN KEY (review_id)
                                      REFERENCES reviews(id)
                                      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS review_votes;
DROP TABLE IF EXISTS reviews;
DROP TYPE IF EXISTS vote_type;
DROP TYPE IF EXISTS review_status;
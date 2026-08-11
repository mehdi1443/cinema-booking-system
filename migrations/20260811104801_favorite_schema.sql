-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE favorites (
                           id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                           user_id BIGINT NOT NULL,
                           movie_id BIGINT,
                           cinema_id BIGINT,
                           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                           CONSTRAINT chk_favorites_target
                               CHECK (
                                   (movie_id IS NOT NULL AND cinema_id IS NULL)
                                       OR
                                   (movie_id IS NULL AND cinema_id IS NOT NULL)
                                   ),

                           CONSTRAINT uq_favorites_movie
                               UNIQUE (user_id, movie_id),

                           CONSTRAINT uq_favorites_cinema
                               UNIQUE (user_id, cinema_id),

                           CONSTRAINT fk_favorites_user
                               FOREIGN KEY (user_id)
                                   REFERENCES users(id)
                                   ON DELETE CASCADE,

                           CONSTRAINT fk_favorites_movie
                               FOREIGN KEY (movie_id)
                                   REFERENCES movies(id)
                                   ON DELETE CASCADE,

                           CONSTRAINT fk_favorites_cinema
                               FOREIGN KEY (cinema_id)
                                   REFERENCES cinemas(id)
                                   ON DELETE CASCADE
);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_movie_id ON favorites(movie_id);
CREATE INDEX idx_favorites_cinema_id ON favorites(cinema_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS favorites CASCADE;
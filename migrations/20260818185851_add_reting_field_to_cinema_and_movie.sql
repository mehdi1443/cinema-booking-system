-- +goose Up

ALTER TABLE movies
    ADD COLUMN rating_avg DECIMAL(3, 2) NOT NULL DEFAULT 0.00,
    ADD COLUMN rating_count INT NOT NULL DEFAULT 0;

ALTER TABLE cinemas
    ADD COLUMN rating_avg DECIMAL(3, 2) NOT NULL DEFAULT 0.00,
    ADD COLUMN rating_count INT NOT NULL DEFAULT 0;


-- +goose Down

ALTER TABLE movies
DROP COLUMN IF EXISTS rating_avg,
    DROP COLUMN IF EXISTS rating_count;

ALTER TABLE cinemas
DROP COLUMN IF EXISTS rating_avg,
    DROP COLUMN IF EXISTS rating_count;
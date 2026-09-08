-- +goose Up
ALTER TABLE cinemas
    ADD COLUMN IF NOT EXISTS owner_user_id BIGINT REFERENCES users(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_cinemas_owner_user_id ON cinemas(owner_user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cinemas_owner_user_id;
ALTER TABLE cinemas DROP COLUMN IF EXISTS owner_user_id;

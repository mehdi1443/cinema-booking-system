-- +goose Up
ALTER TABLE cinemas ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX idx_cinemas_deleted_at ON cinemas(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_cinemas_deleted_at;
ALTER TABLE cinemas DROP COLUMN deleted_at;
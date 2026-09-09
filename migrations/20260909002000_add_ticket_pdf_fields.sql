-- +goose Up
ALTER TABLE tickets
    ADD COLUMN pdf_url TEXT,
    ADD COLUMN pdf_object_key VARCHAR(255),
    ADD COLUMN pdf_generated_at TIMESTAMPTZ,
    ADD CONSTRAINT chk_tickets_pdf_fields CHECK (
        (pdf_url IS NULL AND pdf_object_key IS NULL AND pdf_generated_at IS NULL)
        OR (pdf_url IS NOT NULL AND pdf_object_key IS NOT NULL AND pdf_generated_at IS NOT NULL)
    );

CREATE UNIQUE INDEX uq_tickets_pdf_object_key
    ON tickets(pdf_object_key)
    WHERE pdf_object_key IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS uq_tickets_pdf_object_key;
ALTER TABLE tickets
    DROP CONSTRAINT IF EXISTS chk_tickets_pdf_fields,
    DROP COLUMN IF EXISTS pdf_generated_at,
    DROP COLUMN IF EXISTS pdf_object_key,
    DROP COLUMN IF EXISTS pdf_url;

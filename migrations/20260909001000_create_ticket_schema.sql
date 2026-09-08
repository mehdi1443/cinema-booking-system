-- +goose Up
CREATE TYPE ticket_status AS ENUM ('ACTIVE', 'USED', 'EXPIRED', 'CANCELLED');

CREATE TABLE tickets (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reservation_id BIGINT NOT NULL REFERENCES reservations(id) ON DELETE RESTRICT,
    payment_id BIGINT NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    code VARCHAR(100) NOT NULL,
    status ticket_status NOT NULL DEFAULT 'ACTIVE',
    used_at TIMESTAMPTZ,
    used_by_user_id BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_tickets_reservation UNIQUE (reservation_id),
    CONSTRAINT uq_tickets_payment UNIQUE (payment_id),
    CONSTRAINT uq_tickets_code UNIQUE (code),
    CONSTRAINT chk_tickets_usage_fields CHECK (
        (status = 'USED' AND used_at IS NOT NULL AND used_by_user_id IS NOT NULL)
        OR (status <> 'USED' AND used_at IS NULL AND used_by_user_id IS NULL)
    )
);

CREATE INDEX idx_tickets_user_created ON tickets(user_id, created_at DESC);
CREATE INDEX idx_tickets_status ON tickets(status);

-- +goose Down
DROP TABLE IF EXISTS tickets;
DROP TYPE IF EXISTS ticket_status;

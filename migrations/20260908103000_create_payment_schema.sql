-- +goose Up
CREATE TYPE payment_provider AS ENUM ('MOCK');
CREATE TYPE payment_status AS ENUM ('PENDING', 'SUCCEEDED', 'FAILED');

CREATE TABLE payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reservation_id BIGINT NOT NULL REFERENCES reservations(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL,
    provider payment_provider NOT NULL DEFAULT 'MOCK',
    status payment_status NOT NULL DEFAULT 'PENDING',
    authority VARCHAR(100) NOT NULL,
    reference_id VARCHAR(100),
    idempotency_key VARCHAR(100) NOT NULL,
    failure_reason VARCHAR(255),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_payments_amount_positive CHECK (amount > 0),
    CONSTRAINT uq_payments_authority UNIQUE (authority),
    CONSTRAINT uq_payments_reference_id UNIQUE (reference_id),
    CONSTRAINT uq_payments_user_idempotency UNIQUE (user_id, idempotency_key),
    CONSTRAINT chk_payments_final_fields CHECK (
        (status = 'PENDING' AND paid_at IS NULL AND reference_id IS NULL)
        OR (status = 'SUCCEEDED' AND paid_at IS NOT NULL AND reference_id IS NOT NULL AND failure_reason IS NULL)
        OR (status = 'FAILED' AND paid_at IS NULL AND reference_id IS NULL AND failure_reason IS NOT NULL)
    )
);

CREATE UNIQUE INDEX uq_payments_one_active_per_reservation
    ON payments(reservation_id) WHERE status IN ('PENDING', 'SUCCEEDED');
CREATE INDEX idx_payments_reservation_created ON payments(reservation_id, created_at DESC);
CREATE INDEX idx_payments_user_created ON payments(user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_provider;

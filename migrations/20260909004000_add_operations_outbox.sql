-- +goose Up
CREATE TABLE outbox_events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id BIGINT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'PROCESSING', 'PROCESSED', 'FAILED')),
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_outbox_event UNIQUE (event_type, aggregate_type, aggregate_id)
);

CREATE INDEX idx_outbox_pending ON outbox_events(status, available_at, id)
    WHERE status IN ('PENDING', 'FAILED');

CREATE TABLE audit_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    route VARCHAR(500),
    status_code INT NOT NULL,
    ip_address VARCHAR(64),
    user_agent VARCHAR(500),
    request_id VARCHAR(100),
    duration_ms BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_actor_created ON audit_logs(actor_user_id, created_at DESC);

CREATE TABLE cinema_settlements (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cinema_id BIGINT NOT NULL REFERENCES cinemas(id) ON DELETE RESTRICT,
    period_from TIMESTAMPTZ NOT NULL,
    period_to TIMESTAMPTZ NOT NULL,
    gross_amount BIGINT NOT NULL CHECK (gross_amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PAID')),
    created_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    paid_by_user_id BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (period_to > period_from),
    CHECK ((status = 'PENDING' AND paid_at IS NULL AND paid_by_user_id IS NULL)
        OR (status = 'PAID' AND paid_at IS NOT NULL AND paid_by_user_id IS NOT NULL))
);

CREATE INDEX idx_cinema_settlements_cinema_created
    ON cinema_settlements(cinema_id, created_at DESC);

CREATE TABLE cinema_settlement_items (
    settlement_id BIGINT NOT NULL REFERENCES cinema_settlements(id) ON DELETE RESTRICT,
    payment_id BIGINT NOT NULL UNIQUE REFERENCES payments(id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL CHECK (amount > 0),
    PRIMARY KEY (settlement_id, payment_id)
);

-- +goose Down
DROP TABLE IF EXISTS cinema_settlement_items;
DROP TABLE IF EXISTS cinema_settlements;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS outbox_events;

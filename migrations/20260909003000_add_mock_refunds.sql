-- +goose NO TRANSACTION

-- +goose Up
ALTER TYPE payment_status ADD VALUE IF NOT EXISTS 'REFUNDED';

ALTER TABLE payments
    ADD COLUMN refund_reference_id VARCHAR(100),
    ADD COLUMN refunded_at TIMESTAMPTZ;

ALTER TABLE payments DROP CONSTRAINT chk_payments_final_fields;
ALTER TABLE payments ADD CONSTRAINT chk_payments_final_fields CHECK (
    (status = 'PENDING' AND paid_at IS NULL AND reference_id IS NULL AND failure_reason IS NULL AND refunded_at IS NULL AND refund_reference_id IS NULL)
    OR (status = 'SUCCEEDED' AND paid_at IS NOT NULL AND reference_id IS NOT NULL AND failure_reason IS NULL AND refunded_at IS NULL AND refund_reference_id IS NULL)
    OR (status = 'FAILED' AND paid_at IS NULL AND reference_id IS NULL AND failure_reason IS NOT NULL AND refunded_at IS NULL AND refund_reference_id IS NULL)
    OR (status = 'REFUNDED' AND paid_at IS NOT NULL AND reference_id IS NOT NULL AND failure_reason IS NULL AND refunded_at IS NOT NULL AND refund_reference_id IS NOT NULL)
);

CREATE UNIQUE INDEX uq_payments_refund_reference_id
    ON payments(refund_reference_id)
    WHERE refund_reference_id IS NOT NULL;

-- +goose Down
ALTER TABLE payments DROP CONSTRAINT chk_payments_final_fields;
DROP INDEX IF EXISTS uq_payments_refund_reference_id;
DROP INDEX IF EXISTS uq_payments_one_active_per_reservation;
UPDATE payments SET status = 'SUCCEEDED' WHERE status = 'REFUNDED';
ALTER TABLE payments
    DROP COLUMN IF EXISTS refunded_at,
    DROP COLUMN IF EXISTS refund_reference_id;

ALTER TABLE payments ALTER COLUMN status DROP DEFAULT;
ALTER TABLE payments ALTER COLUMN status TYPE VARCHAR(20) USING status::TEXT;
DROP TYPE payment_status;
CREATE TYPE payment_status AS ENUM ('PENDING', 'SUCCEEDED', 'FAILED');
ALTER TABLE payments ALTER COLUMN status TYPE payment_status USING status::payment_status;
ALTER TABLE payments ALTER COLUMN status SET DEFAULT 'PENDING';

CREATE UNIQUE INDEX uq_payments_one_active_per_reservation
    ON payments(reservation_id) WHERE status IN ('PENDING', 'SUCCEEDED');

ALTER TABLE payments ADD CONSTRAINT chk_payments_final_fields CHECK (
    (status = 'PENDING' AND paid_at IS NULL AND reference_id IS NULL)
    OR (status = 'SUCCEEDED' AND paid_at IS NOT NULL AND reference_id IS NOT NULL AND failure_reason IS NULL)
    OR (status = 'FAILED' AND paid_at IS NULL AND reference_id IS NULL AND failure_reason IS NOT NULL)
);

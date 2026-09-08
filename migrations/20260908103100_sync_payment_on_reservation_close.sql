-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fail_pending_payments_for_closed_reservation()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IN ('EXPIRED'::reservation_status, 'CANCELLED'::reservation_status)
       AND NEW.status IS DISTINCT FROM OLD.status THEN
        UPDATE payments
           SET status = 'FAILED',
               failure_reason = CASE
                   WHEN NEW.status = 'EXPIRED'::reservation_status THEN 'RESERVATION_EXPIRED'
                   ELSE 'RESERVATION_CANCELLED'
               END,
               updated_at = NOW()
         WHERE reservation_id = NEW.id
           AND status = 'PENDING';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_fail_payment_on_reservation_close
AFTER UPDATE OF status ON reservations
FOR EACH ROW
EXECUTE FUNCTION fail_pending_payments_for_closed_reservation();

-- +goose Down
DROP TRIGGER IF EXISTS trg_fail_payment_on_reservation_close ON reservations;
DROP FUNCTION IF EXISTS fail_pending_payments_for_closed_reservation();

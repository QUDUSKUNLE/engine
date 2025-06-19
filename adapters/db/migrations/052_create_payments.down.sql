-- Drop trigger first
DROP TRIGGER IF EXISTS update_payment_timestamp ON payments;

-- Drop indexes
DROP INDEX IF EXISTS idx_payments_appointment;
DROP INDEX IF EXISTS idx_payments_patient;
DROP INDEX IF EXISTS idx_payments_centre;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_date;
DROP INDEX IF EXISTS idx_payments_created;

-- Drop table
DROP TABLE IF EXISTS payments;

-- Drop type
DROP TYPE IF EXISTS payment_method;

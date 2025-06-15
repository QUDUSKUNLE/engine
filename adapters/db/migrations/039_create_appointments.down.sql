-- Drop triggers first
DROP TRIGGER IF EXISTS update_appointment_timestamp ON appointments;
DROP TRIGGER IF EXISTS validate_appointment ON appointments;

-- Drop functions
DROP FUNCTION IF EXISTS validate_appointment_time_slot() CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_patient;
DROP INDEX IF EXISTS idx_appointments_schedule;
DROP INDEX IF EXISTS idx_appointments_centre;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_date;
DROP INDEX IF EXISTS idx_appointments_payment_status;
DROP INDEX IF EXISTS idx_appointments_created;

-- Drop table
DROP TABLE IF EXISTS appointments CASCADE;

-- Drop ENUMs
DROP TYPE IF EXISTS appointment_status CASCADE;
DROP TYPE IF EXISTS payment_status CASCADE;

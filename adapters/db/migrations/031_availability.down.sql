-- Drop triggers first
DROP TRIGGER IF EXISTS validate_test_type_trigger ON diagnostic_schedules;
DROP TRIGGER IF EXISTS check_slot_availability_trigger ON diagnostic_schedules;
DROP TRIGGER IF EXISTS update_diagnostic_centre_availability_updated_at ON diagnostic_centre_availability;

-- Drop functions
DROP FUNCTION IF EXISTS validate_test_type() CASCADE;
DROP FUNCTION IF EXISTS check_slot_availability() CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS idx_diagnostic_centre_availability_centre;
DROP INDEX IF EXISTS idx_diagnostic_centre_availability_day_time;

-- Drop table
DROP TABLE IF EXISTS diagnostic_centre_availability CASCADE;

-- Drop type
DROP TYPE IF EXISTS weekday CASCADE;

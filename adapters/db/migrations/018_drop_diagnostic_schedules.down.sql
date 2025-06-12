-- Drop indexes
DROP INDEX IF EXISTS idx_diagnostic_schedules_id_user_id;
DROP INDEX IF EXISTS idx_diagnostic_schedules_centre_date;
DROP INDEX IF EXISTS idx_diagnostic_schedules_centre_status_time;

-- Drop ENUM types
DROP TYPE IF EXISTS schedule_acceptance_status CASCADE;
DROP TYPE IF EXISTS schedule_status CASCADE;
DROP TYPE IF EXISTS test_type CASCADE;

-- Drop table
DROP TABLE IF EXISTS diagnostic_schedules CASCADE;

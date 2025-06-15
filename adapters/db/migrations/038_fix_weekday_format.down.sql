-- Drop the function so it can be recreated by the up migration
DROP FUNCTION IF EXISTS check_slot_availability() CASCADE;

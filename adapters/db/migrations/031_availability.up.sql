-- Drop any existing objects that might conflict
DO $$ 
BEGIN
    -- Drop triggers if they exist
    DROP TRIGGER IF EXISTS validate_test_type_trigger ON diagnostic_schedules;
    DROP TRIGGER IF EXISTS check_slot_availability_trigger ON diagnostic_schedules;
    DROP TRIGGER IF EXISTS update_diagnostic_centre_availability_updated_at ON diagnostic_centre_availability;
    
    -- Drop functions if they exist
    DROP FUNCTION IF EXISTS validate_test_type() CASCADE;
    DROP FUNCTION IF EXISTS check_slot_availability() CASCADE;
    
    -- Drop table if exists
    DROP TABLE IF EXISTS diagnostic_centre_availability CASCADE;
    
    -- Drop type if exists
    DROP TYPE IF EXISTS weekday CASCADE;
END $$;

-- Create weekday enum for availability
CREATE TYPE weekday AS ENUM ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday');

-- Create the diagnostic centre availability table
CREATE TABLE IF NOT EXISTS diagnostic_centre_availability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    diagnostic_centre_id UUID NOT NULL REFERENCES diagnostic_centres(id) ON DELETE CASCADE,
    day_of_week weekday NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_appointments INTEGER DEFAULT 0,
    slot_duration INTERVAL NOT NULL DEFAULT '30 minutes',
    break_time INTERVAL DEFAULT '5 minutes',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_time_range CHECK (start_time < end_time),
    CONSTRAINT check_max_appointments CHECK (max_appointments >= 0),
    CONSTRAINT check_slot_duration CHECK (slot_duration > '0 minutes'),
    UNIQUE (diagnostic_centre_id, day_of_week)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_diagnostic_centre_availability_centre 
    ON diagnostic_centre_availability(diagnostic_centre_id);
CREATE INDEX IF NOT EXISTS idx_diagnostic_centre_availability_day_time 
    ON diagnostic_centre_availability(day_of_week, start_time, end_time);

-- Add trigger to update updated_at timestamp
DROP TRIGGER IF EXISTS update_diagnostic_centre_availability_updated_at ON diagnostic_centre_availability;
CREATE TRIGGER update_diagnostic_centre_availability_updated_at
    BEFORE UPDATE ON diagnostic_centre_availability
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create a function to validate test types against diagnostic center's available tests
CREATE OR REPLACE FUNCTION validate_test_type() RETURNS TRIGGER AS $$
BEGIN
    -- Check if the test type is available at the diagnostic center
    IF NOT EXISTS (
        SELECT 1 FROM diagnostic_centres
        WHERE id = NEW.diagnostic_centre_id
        AND NEW.test_type = ANY(available_tests)
    ) THEN
        RAISE EXCEPTION 'Test type % is not available at this diagnostic center', NEW.test_type
            USING ERRCODE = 'check_violation';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger to validate test type on insert and update
DROP TRIGGER IF EXISTS validate_test_type_trigger ON diagnostic_schedules;
CREATE TRIGGER validate_test_type_trigger 
    BEFORE INSERT OR UPDATE OF test_type ON diagnostic_schedules
    FOR EACH ROW EXECUTE FUNCTION validate_test_type();

-- Create a function to check slot availability when scheduling
CREATE OR REPLACE FUNCTION check_slot_availability() RETURNS TRIGGER AS $$
DECLARE
    slot_count INTEGER;
    max_slots INTEGER;
    slot_range tsrange;
BEGIN
    -- Get the day of week for the requested schedule time
    -- Get availability settings for that day
    SELECT 
        a.max_appointments,
        tsrange(
            date_trunc('day', NEW.schedule_time) + a.start_time,
            date_trunc('day', NEW.schedule_time) + a.end_time
        ) INTO max_slots, slot_range
    FROM diagnostic_centre_availability a
    WHERE a.diagnostic_centre_id = NEW.diagnostic_centre_id
    AND a.day_of_week = LOWER(TO_CHAR(NEW.schedule_time, 'day'));

    IF max_slots IS NULL THEN
        RAISE EXCEPTION 'No availability set for this day'
            USING ERRCODE = 'check_violation';
    END IF;

    -- Check if schedule_time falls within the available time range
    IF NOT (NEW.schedule_time <@ slot_range) THEN
        RAISE EXCEPTION 'Schedule time is outside available hours'
            USING ERRCODE = 'check_violation';
    END IF;

    -- Count existing appointments in the same time slot
    SELECT COUNT(*) INTO slot_count
    FROM diagnostic_schedules s
    WHERE s.diagnostic_centre_id = NEW.diagnostic_centre_id
    AND s.schedule_time = NEW.schedule_time
    AND s.id != NEW.id; -- Exclude current record for updates

    IF slot_count >= max_slots THEN
        RAISE EXCEPTION 'No available slots at this time'
            USING ERRCODE = 'check_violation';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger to validate slot availability on insert and update
DROP TRIGGER IF EXISTS check_slot_availability_trigger ON diagnostic_schedules;
CREATE TRIGGER check_slot_availability_trigger
    BEFORE INSERT OR UPDATE OF schedule_time ON diagnostic_schedules
    FOR EACH ROW EXECUTE FUNCTION check_slot_availability();

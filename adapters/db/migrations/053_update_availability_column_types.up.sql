-- First, drop constraints and default values
ALTER TABLE diagnostic_centre_availability
    DROP CONSTRAINT IF EXISTS check_slot_duration,
    DROP CONSTRAINT IF EXISTS check_break_time,
    DROP CONSTRAINT IF EXISTS check_day_of_week,
    ALTER COLUMN slot_duration DROP DEFAULT,
    ALTER COLUMN break_time DROP DEFAULT;

-- Convert day_of_week from weekday enum to text
ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN day_of_week TYPE TEXT;

-- Convert slot_duration from interval to integer (minutes)
ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN slot_duration TYPE INTEGER 
    USING EXTRACT(EPOCH FROM slot_duration)/60,
    ALTER COLUMN slot_duration SET DEFAULT 30;

-- Convert break_time from interval to integer (minutes)
ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN break_time TYPE INTEGER 
    USING EXTRACT(EPOCH FROM break_time)/60,
    ALTER COLUMN break_time SET DEFAULT 5;

-- Update constraints
ALTER TABLE diagnostic_centre_availability
    DROP CONSTRAINT IF EXISTS check_slot_duration;

ALTER TABLE diagnostic_centre_availability
    ADD CONSTRAINT check_slot_duration CHECK (slot_duration > 0),
    ADD CONSTRAINT check_break_time CHECK (break_time >= 0),
    ADD CONSTRAINT check_day_of_week CHECK (day_of_week IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'));

-- Update the check_slot_availability function to work with the new types
CREATE OR REPLACE FUNCTION check_slot_availability() RETURNS TRIGGER AS $$
DECLARE
    slot_count INTEGER;
    max_slots INTEGER;
    slot_range tstzrange;
    day_name text;
BEGIN
    -- Convert the day name to lowercase and trim spaces
    day_name := TRIM(LOWER(TO_CHAR(NEW.schedule_time AT TIME ZONE 'UTC', 'day')));

    -- Get availability settings for that day
    SELECT 
        a.max_appointments,
        tstzrange(
            date_trunc('day', NEW.schedule_time AT TIME ZONE 'UTC') + a.start_time,
            date_trunc('day', NEW.schedule_time AT TIME ZONE 'UTC') + a.end_time
        ) INTO max_slots, slot_range
    FROM diagnostic_centre_availability a
    WHERE a.diagnostic_centre_id = NEW.diagnostic_centre_id
    AND a.day_of_week = day_name;

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
    AND s.id != NEW.id;

    IF slot_count >= max_slots THEN
        RAISE EXCEPTION 'No available slots at this time'
            USING ERRCODE = 'check_violation';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

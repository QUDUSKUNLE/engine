-- Drop new constraints first
ALTER TABLE diagnostic_centre_availability
    DROP CONSTRAINT IF EXISTS check_slot_duration,
    DROP CONSTRAINT IF EXISTS check_break_time,
    DROP CONSTRAINT IF EXISTS check_day_of_week;

-- First drop the default values
ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN slot_duration DROP DEFAULT,
    ALTER COLUMN break_time DROP DEFAULT;

-- Convert columns back to their original types
ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN day_of_week TYPE weekday 
    USING day_of_week::weekday;

ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN slot_duration TYPE INTERVAL 
    USING (slot_duration || ' minutes')::interval,
    ALTER COLUMN slot_duration SET DEFAULT '30 minutes';

ALTER TABLE diagnostic_centre_availability
    ALTER COLUMN break_time TYPE INTERVAL 
    USING (break_time || ' minutes')::interval,
    ALTER COLUMN break_time SET DEFAULT '5 minutes';

-- Restore original constraint
ALTER TABLE diagnostic_centre_availability
    ADD CONSTRAINT check_slot_duration CHECK (slot_duration > '0 minutes');

-- Restore original slot availability function
CREATE OR REPLACE FUNCTION check_slot_availability() RETURNS TRIGGER AS $$
DECLARE
    slot_count INTEGER;
    max_slots INTEGER;
    slot_range tstzrange;
    day_name weekday;
BEGIN
    -- Convert the day name to lowercase, trim spaces, and cast to weekday enum
    day_name := TRIM(LOWER(TO_CHAR(NEW.schedule_time AT TIME ZONE 'UTC', 'day')))::weekday;

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

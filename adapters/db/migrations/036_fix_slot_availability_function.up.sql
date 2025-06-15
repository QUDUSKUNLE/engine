-- Update the check_slot_availability function to use tstzrange
CREATE OR REPLACE FUNCTION check_slot_availability() RETURNS TRIGGER AS $$
DECLARE
    slot_count INTEGER;
    max_slots INTEGER;
    slot_range tstzrange;
BEGIN
    -- Get the day of week for the requested schedule time
    -- Get availability settings for that day
    SELECT 
        a.max_appointments,
        tstzrange(
            date_trunc('day', NEW.schedule_time AT TIME ZONE 'UTC') + a.start_time,
            date_trunc('day', NEW.schedule_time AT TIME ZONE 'UTC') + a.end_time
        ) INTO max_slots, slot_range
    FROM diagnostic_centre_availability a
    WHERE a.diagnostic_centre_id = NEW.diagnostic_centre_id
    AND a.day_of_week = LOWER(TO_CHAR(NEW.schedule_time AT TIME ZONE 'UTC', 'day'));

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

-- Create payment status ENUM
CREATE TYPE payment_status AS ENUM (
    'pending',
    'success',
    'failed',
    'refunded',
    'cancelled'
);

-- Create appointment status ENUM
CREATE TYPE appointment_status AS ENUM (
    'pending',
    'confirmed',
    'in_progress',
    'completed',
    'cancelled',
    'rescheduled'
);

-- Create appointments table
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign keys
    patient_id UUID NOT NULL REFERENCES users(id),
    schedule_id UUID NOT NULL REFERENCES diagnostic_schedules(id),
    diagnostic_centre_id UUID NOT NULL REFERENCES diagnostic_centres(id),
    
    -- Appointment details
    appointment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    time_slot VARCHAR(20) NOT NULL, -- e.g. "09:00-09:30"
    status appointment_status NOT NULL DEFAULT 'pending',
    
    -- Payment info
    payment_id UUID NULL, -- Will be filled when payment is initiated
    payment_status payment_status NULL,
    payment_amount DECIMAL(10,2) NULL,
    payment_date TIMESTAMP WITH TIME ZONE NULL,
    
    -- Tracking
    check_in_time TIMESTAMP WITH TIME ZONE NULL,
    completion_time TIMESTAMP WITH TIME ZONE NULL,
    notes TEXT NULL,

    -- Cancellation/Rescheduling
    cancellation_reason TEXT NULL,
    cancelled_by UUID NULL REFERENCES users(id),
    cancellation_time TIMESTAMP WITH TIME ZONE NULL,
    cancellation_fee DECIMAL(10,2) NULL,
    
    original_appointment_id UUID NULL REFERENCES appointments(id), -- For rescheduled appointments
    rescheduling_reason TEXT NULL,
    rescheduled_by UUID NULL REFERENCES users(id),
    rescheduling_time TIMESTAMP WITH TIME ZONE NULL,
    rescheduling_fee DECIMAL(10,2) NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CHECK ((status = 'cancelled' AND cancellation_reason IS NOT NULL AND cancelled_by IS NOT NULL AND cancellation_time IS NOT NULL) OR 
           (status != 'cancelled' AND cancellation_reason IS NULL AND cancelled_by IS NULL AND cancellation_time IS NULL)),
    CHECK ((status = 'rescheduled' AND rescheduling_reason IS NOT NULL AND rescheduled_by IS NOT NULL AND rescheduling_time IS NOT NULL AND original_appointment_id IS NOT NULL) OR
           (status != 'rescheduled' AND rescheduling_reason IS NULL AND rescheduled_by IS NULL AND rescheduling_time IS NULL AND original_appointment_id IS NULL)),
    CHECK (payment_date IS NULL OR payment_date <= CURRENT_TIMESTAMP),
    CHECK (check_in_time IS NULL OR check_in_time <= CURRENT_TIMESTAMP),
    CHECK (completion_time IS NULL OR completion_time <= CURRENT_TIMESTAMP),
    CHECK (cancellation_time IS NULL OR cancellation_time <= CURRENT_TIMESTAMP),
    CHECK (rescheduling_time IS NULL OR rescheduling_time <= CURRENT_TIMESTAMP),
    CHECK (check_in_time IS NULL OR check_in_time >= appointment_date),
    CHECK (completion_time IS NULL OR completion_time >= check_in_time),
    CHECK (payment_amount >= 0),
    CHECK (cancellation_fee >= 0),
    CHECK (rescheduling_fee >= 0)
);

-- Add indexes for common queries
CREATE INDEX idx_appointments_patient ON appointments(patient_id);
CREATE INDEX idx_appointments_schedule ON appointments(schedule_id);
CREATE INDEX idx_appointments_centre ON appointments(diagnostic_centre_id);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_date ON appointments(appointment_date);
CREATE INDEX idx_appointments_payment_status ON appointments(payment_status);
CREATE INDEX idx_appointments_created ON appointments(created_at DESC);

-- Function to update updated_at timestamp
CREATE TRIGGER update_appointment_timestamp BEFORE UPDATE
ON appointments FOR EACH ROW EXECUTE FUNCTION
update_updated_at_column();

-- Function to validate appointment time slots
CREATE OR REPLACE FUNCTION validate_appointment_time_slot() 
RETURNS TRIGGER AS $$
BEGIN
    -- Make sure appointment isn't in the past
    IF NEW.appointment_date < CURRENT_TIMESTAMP THEN
        RAISE EXCEPTION 'Cannot create appointments in the past';
    END IF;

    -- Validate time slot format (HH:MM-HH:MM)
    IF NOT NEW.time_slot ~ '^([0-1][0-9]|2[0-3]):[0-5][0-9]-([0-1][0-9]|2[0-3]):[0-5][0-9]$' THEN
        RAISE EXCEPTION 'Invalid time slot format. Must be HH:MM-HH:MM';
    END IF;

    -- Check if the schedule exists and is ACCEPTED
    IF NOT EXISTS (
        SELECT 1 FROM diagnostic_schedules 
        WHERE id = NEW.schedule_id 
        AND diagnostic_centre_id = NEW.diagnostic_centre_id
        AND acceptance_status = 'ACCEPTED'
    ) THEN
        RAISE EXCEPTION 'Invalid or unconfirmed schedule';
    END IF;

    -- Check for double booking in the same time slot
    IF EXISTS (
        SELECT 1 FROM appointments
        WHERE diagnostic_centre_id = NEW.diagnostic_centre_id
        AND appointment_date = NEW.appointment_date
        AND time_slot = NEW.time_slot
        AND status NOT IN ('cancelled', 'rescheduled')
        AND id != NEW.id
    ) THEN
        RAISE EXCEPTION 'Time slot already booked';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger for appointment validation
CREATE TRIGGER validate_appointment
    BEFORE INSERT OR UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION validate_appointment_time_slot();

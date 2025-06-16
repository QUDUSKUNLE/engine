-- Add reminder tracking to appointments
ALTER TABLE appointments
ADD COLUMN reminder_sent boolean DEFAULT false,
ADD COLUMN reminder_sent_at timestamp with time zone;

-- Add index for reminder queries
CREATE INDEX idx_appointments_reminder ON appointments(reminder_sent, appointment_date)
WHERE status = 'confirmed' AND reminder_sent = false;

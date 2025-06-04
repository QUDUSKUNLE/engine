-- Add new values to schedule_type enum, one at a time
ALTER TYPE user_enum ADD VALUE IF NOT EXISTS 'DIAGNOSTIC_MANAGER';

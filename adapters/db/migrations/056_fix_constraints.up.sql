-- Drop existing constraint if it exists
ALTER TABLE IF EXISTS diagnostic_centres DROP CONSTRAINT IF EXISTS unique_admin_id;

-- Add the constraint back
ALTER TABLE diagnostic_centres ADD CONSTRAINT unique_admin_id UNIQUE (admin_id);


-- Add the new unique constraint
ALTER TABLE diagnostic_centres ADD CONSTRAINT unique_admin_id UNIQUE (admin_id);

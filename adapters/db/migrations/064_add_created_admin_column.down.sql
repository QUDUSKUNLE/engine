-- Remove foreign key constraint first
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_created_admin;

-- Remove created_admin column
ALTER TABLE users
DROP COLUMN IF EXISTS created_admin;

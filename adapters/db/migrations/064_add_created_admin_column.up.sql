-- Add created_admin column to users table
ALTER TABLE users
ADD COLUMN created_admin UUID NULL;

-- Add foreign key constraint referencing users table
ALTER TABLE users
ADD CONSTRAINT fk_created_admin
FOREIGN KEY (created_admin)
REFERENCES users(id)
ON DELETE SET NULL;

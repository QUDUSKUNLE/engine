ALTER TABLE users
    DROP COLUMN IF EXISTS phone_number;

DROP INDEX IF EXISTS idx_users_phone;

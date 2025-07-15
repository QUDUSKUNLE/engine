
-- Update all users to PATIENT
UPDATE users
SET user_type = 'PATIENT'
WHERE user_type = 'USER';

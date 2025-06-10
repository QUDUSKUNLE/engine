DROP TABLE IF EXISTS diagnostic_centres CASCADE;
DROP TYPE IF EXISTS doctor;
DROP TYPE IF EXISTS available_tests;

-- DROP Constraints
ALTER TABLE IF EXISTS diagnostic_centres DROP CONSTRAINT IF EXISTS diagnostic_centres_pkey;
ALTER TABLE IF EXISTS diagnostic_centres DROP CONSTRAINT IF EXISTS diagnostic_centres_admin_id_fkey;
ALTER TABLE IF EXISTS diagnostic_centres DROP CONSTRAINT IF EXISTS diagnostic_centres_created_by_fkey;

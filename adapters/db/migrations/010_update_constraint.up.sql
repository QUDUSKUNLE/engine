-- Drop the old unique constraint if it exists (replace constraint name if different)
ALTER TABLE diagnostic_centres
DROP CONSTRAINT IF EXISTS diagnostic_centres_diagnostic_centre_name_created_by_longitude_latitude_address_key;

-- Add the new unique constraint
ALTER TABLE diagnostic_centres
ADD CONSTRAINT diagnostic_centres_latitude_longitude_key UNIQUE (latitude, longitude);

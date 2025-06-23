-- Change doctors column to text array
ALTER TABLE diagnostic_centres 
ALTER COLUMN doctors TYPE text[] 
USING ARRAY[doctors]::text[];

-- Change available_tests column to text array
ALTER TABLE diagnostic_centres 
ALTER COLUMN available_tests TYPE text[]
USING ARRAY[available_tests]::text[];

-- Update existing data to have proper casing
UPDATE diagnostic_centres
SET doctors = ARRAY(
    SELECT CASE WHEN lower(unnest) = 'female' THEN 'Female'
                WHEN lower(unnest) = 'male' THEN 'Male'
                ELSE unnest
           END
    FROM unnest(doctors)
);

-- Update existing data to proper format
UPDATE diagnostic_centres
SET available_tests = ARRAY(
    SELECT upper(replace(unnest, ' ', '_'))
    FROM unnest(available_tests)
);

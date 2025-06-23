-- Change doctors column back to Doctor array
ALTER TABLE diagnostic_centres 
ALTER COLUMN doctors TYPE Doctor[] 
USING ARRAY(
    SELECT unnest::Doctor
    FROM unnest(doctors)
);

-- Change available_tests column back to AvailableTests array
ALTER TABLE diagnostic_centres 
ALTER COLUMN available_tests TYPE AvailableTests[]
USING ARRAY(
    SELECT unnest::AvailableTests
    FROM unnest(available_tests)
);

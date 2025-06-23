-- Drop the existing trigger since it depends on the enum
DROP TRIGGER IF EXISTS validate_test_type_trigger ON diagnostic_schedules;

-- Drop the check constraint if it exists
ALTER TABLE diagnostic_schedules DROP CONSTRAINT IF EXISTS valid_test_type_check;

-- First remove the default value that depends on the enum
ALTER TABLE diagnostic_schedules ALTER COLUMN test_type DROP DEFAULT;

-- Convert the column to text while preserving the data
ALTER TABLE diagnostic_schedules 
    ALTER COLUMN test_type TYPE text USING test_type::text;

-- Drop the enum type since we no longer need it
DROP TYPE IF EXISTS test_type;

-- Add a new default value using text type
ALTER TABLE diagnostic_schedules ALTER COLUMN test_type SET DEFAULT 'OTHER';

-- Add the check constraint to maintain data validation
ALTER TABLE diagnostic_schedules 
ADD CONSTRAINT valid_test_type_check CHECK (
    test_type IN (
        'BLOOD_TEST', 'URINE_TEST', 'X_RAY', 'MRI', 'CT_SCAN', 'ULTRASOUND', 'ECG',
        'EEG', 'BIOPSY', 'SKIN_TEST', 'IMMUNOLOGY_TEST', 'HORMONE_TEST', 'VIRAL_TEST',
        'BACTERIAL_TEST', 'PARASITIC_TEST', 'FUNGAL_TEST', 'MOLECULAR_TEST', 'TOXICOLOGY_TEST',
        'ECHO', 'COVID_19_TEST', 'BLOOD_SUGAR_TEST', 'LIPID_PROFILE', 'HEMOGLOBIN_TEST',
        'THYROID_TEST', 'LIVER_FUNCTION_TEST', 'KIDNEY_FUNCTION_TEST', 'URIC_ACID_TEST',
        'VITAMIN_D_TEST', 'VITAMIN_B12_TEST', 'HEMOGRAM', 'COMPLETE_BLOOD_COUNT',
        'BLOOD_GROUPING', 'HEPATITIS_B_TEST', 'HEPATITIS_C_TEST', 'HIV_TEST',
        'MALARIA_TEST', 'DENGUE_TEST', 'TYPHOID_TEST', 'COVID_19_ANTIBODY_TEST',
        'COVID_19_RAPID_ANTIGEN_TEST', 'COVID_19_RT_PCR_TEST', 'PREGNANCY_TEST',
        'ALLERGY_TEST', 'GENETIC_TEST', 'OTHER'
    )
);

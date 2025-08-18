ALTER TABLE medical_records DROP CONSTRAINT medical_records_uploader_type_check;

ALTER TABLE medical_records
ADD CONSTRAINT medical_records_uploader_type_check
CHECK (
  uploader_type = ANY(
    ARRAY['DIAGNOSTIC_CENTRE_MANAGER', 'DIAGNOSTIC_CENTRE_OWNER']::user_enum[]
  )
);

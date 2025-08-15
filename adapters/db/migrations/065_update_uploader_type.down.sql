-- Down
ALTER TABLE medical_records
DROP CONSTRAINT IF EXISTS medical_records_uploader_type_check;

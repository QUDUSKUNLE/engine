CREATE TYPE document_type AS ENUM (
  'LAB_REPORT',
  'PRESCRIPTION',
  'DISCHARGE_SUMMARY',
  'IMAGING',
  'VACCINATION',
  'ALLERGY',
  'SURGERY',
  'CHRONIC_CONDITION',
  'FAMILY_HISTORY'
);

CREATE TABLE medical_records (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  -- User identity (who the record belongs to)
  user_id UUID NOT NULL REFERENCES users(id),

  -- Uploader identity (who uploaded it)
  uploader_id UUID NOT NULL,
  uploader_type user_enum NOT NULL CHECK (uploader_type NOT IN ('USER', 'ADMIN', 'HOSPITAL')),

  -- Schedule ID (record linked to a diagnostic schedule)
  schedule_id UUID NOT NULL REFERENCES diagnostic_schedules(id),

  -- Record details
  title VARCHAR(255) NOT NULL,
  document_type document_type NOT NULL,
  file_path TEXT NOT NULL,
  file_type VARCHAR(50), -- e.g., pdf, jpg, dicom
  document_date DATE,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  provider_name VARCHAR(255),
  specialty VARCHAR(100),
  is_shared BOOLEAN DEFAULT FALSE,
  shared_until TIMESTAMP,

  -- Metadata
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Trigger function to validate uploader_id based on uploader_type
CREATE OR REPLACE FUNCTION validate_uploader_id()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.uploader_type = 'DOCTOR' THEN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = NEW.uploader_id) THEN
      RAISE EXCEPTION 'uploader_id % does not exist in users table', NEW.uploader_id;
    END IF;
  ELSIF NEW.uploader_type = 'DIAGNOSTIC_CENTRE_MANAGER' THEN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = NEW.uploader_id) THEN
      RAISE EXCEPTION 'uploader_id % does not exist in users table', NEW.uploader_id;
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for INSERT and UPDATE on medical_records
CREATE TRIGGER trg_validate_uploader_id
BEFORE INSERT OR UPDATE ON medical_records
FOR EACH ROW EXECUTE FUNCTION validate_uploader_id();

-- CREATE TABLE diagnostic_availability (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     diagnostic_centres UUID REFERENCES diagnostic_centres(id),
--     day_of_week ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'),
--     start_time TIME,
--     end_time TIME,
--     max_appointments INTEGER DEFAULT 0,
--     test_type VARCHAR(255) -- nullable if availability is generic
-- );

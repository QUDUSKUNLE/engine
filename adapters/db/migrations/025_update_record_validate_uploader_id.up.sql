-- Migration: Update validate_uploader_id() to check diagnostic_centres for DIAGNOSTIC_CENTRE_MANAGER
-- Created: 2025-06-12

-- Update the validate_uploader_id function to check diagnostic_centres for DIAGNOSTIC_CENTRE_MANAGER
CREATE OR REPLACE FUNCTION validate_uploader_id()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.uploader_type = 'DIAGNOSTIC_CENTRE_MANAGER' THEN
    IF NOT EXISTS (SELECT 1 FROM diagnostic_centres WHERE id = NEW.uploader_id AND admin_id = NEW.uploader_admin_id) THEN
      RAISE EXCEPTION 'uploader_id % or uploader_admin_id % do not exist in diagnostic_centres table', NEW.uploader_id, NEW.uploader_admin_id;
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

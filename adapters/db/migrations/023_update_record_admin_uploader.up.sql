ALTER TABLE medical_records
ADD COLUMN uploader_admin_id UUID REFERENCES users(id);

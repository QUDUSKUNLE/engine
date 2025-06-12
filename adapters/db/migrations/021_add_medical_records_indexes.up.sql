-- Migration: Add indexes for medical_records queries optimization
-- Created: 2025-06-12

-- Index for queries filtering by user_id
CREATE INDEX IF NOT EXISTS idx_medical_records_user_id ON medical_records(user_id);

-- Index for queries filtering by uploader_id
CREATE INDEX IF NOT EXISTS idx_medical_records_uploader_id ON medical_records(uploader_id);

-- Composite index for user_id and ordering by created_at
CREATE INDEX IF NOT EXISTS idx_medical_records_user_id_created_at ON medical_records(user_id, created_at DESC);

-- Composite index for uploader_id and ordering by created_at
CREATE INDEX IF NOT EXISTS idx_medical_records_uploader_id_created_at ON medical_records(uploader_id, created_at DESC);

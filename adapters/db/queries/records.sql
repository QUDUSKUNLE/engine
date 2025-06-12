-- Create a Medical Record
-- name: CreateMedicalRecord :one
INSERT INTO medical_records (
  user_id,
  uploader_id,
  uploader_admin_id,
  uploader_type,
  schedule_id,
  title,
  document_type,
  document_date,
  file_path,
  file_type,
  uploaded_at,
  provider_name,
  specialty,
  is_shared,
  shared_until
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- Get a Medical Record
-- name: GetMedicalRecord :one
SELECT id, user_id, uploader_id, uploader_admin_id, uploader_type, schedule_id, title, document_type, document_date, file_path, file_type, uploaded_at, provider_name, specialty, is_shared, shared_until, created_at, updated_at
FROM medical_records WHERE id = $1 AND user_id = $2;

-- Get an Uploader Medical Record
-- Retrieves a medical record by its ID and uploader ID.
-- name: GetUploaderMedicalRecord :one
SELECT id, user_id, uploader_id, uploader_admin_id, uploader_type, schedule_id, title, document_type, document_date, file_path, file_type, uploaded_at, provider_name, specialty, is_shared, shared_until, created_at, updated_at
FROM medical_records WHERE id = $1 AND uploader_id = $2 AND uploader_admin_id = $3;

-- Get Medical Records
-- Retrieves a paginated list of medical records for a specific user, ordered by creation date (most recent first).
-- name: GetMedicalRecords :many
SELECT id, user_id, uploader_id, uploader_admin_id, uploader_type, schedule_id, title, document_type, document_date, file_path, file_type, uploaded_at, provider_name, specialty, is_shared, shared_until, created_at FROM medical_records
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Get Uploader Medical Records
-- Retrieves a paginated list of medical records uploaded by a specific uploader, explicitly ordered by creation date in descending order (most recent first).
-- name: GetUploaderMedicalRecords :many
SELECT id, user_id, uploader_id, uploader_admin_id, uploader_type, schedule_id, title, document_type, document_date, file_path, file_type, uploaded_at, provider_name, specialty, is_shared, shared_until, created_at FROM medical_records
WHERE uploader_id = $1 AND uploader_admin_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- Uploader Update a Medical Record
-- Updates a medical record by uploader, allowing partial updates to fields. Only the uploader can update their own records. Updates the 'updated_at' timestamp.
-- name: UpdateMedicalRecordByUploader :one
UPDATE medical_records
SET 
  title = COALESCE($1, title),
  document_type = COALESCE($2, document_type),
  document_date = COALESCE($3, document_date),
  file_path = COALESCE($4, file_path),
  file_type = COALESCE($5, file_type),
  uploaded_at = COALESCE($6, uploaded_at),
  provider_name = COALESCE($7, provider_name),
  specialty = COALESCE($8, specialty),
  is_shared = COALESCE($9, is_shared),
  shared_until = COALESCE($10, shared_until),
  updated_at = NOW()
WHERE id = $11 AND uploader_id = $12 AND uploader_admin_id = $13
RETURNING *;



-- Create a diagnostic schedule
-- name: Create_Diagnostic_Schedule :one
INSERT INTO diagnostic_schedules (
  user_id,
  diagnostic_centre_id,
  schedule_time,
  test_type,
  schedule_status,
  doctor,
  acceptance_status,
  notes
) VALUES (
   $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- Get Diagnostic Schedule
-- name: Get_Diagnostic_Schedule :one
SELECT * FROM diagnostic_schedules WHERE id = $1 AND user_id = $2;

-- Get Diagnostic Schedules
-- name: Get_Diagnostic_Schedules :many
SELECT * FROM diagnostic_schedules
WHERE user_id = $1
ORDER BY schedule_time DESC
LIMIT $2 OFFSET $3;

-- Update a diagnostic schedule
-- name: Update_Diagnostic_Schedule :one
UPDATE diagnostic_schedules
SET
  schedule_time = COALESCE($1, schedule_time),
  test_type = COALESCE($2, test_type),
  schedule_status = COALESCE($3, schedule_status),
  notes = COALESCE($4, notes),
  doctor = COALESCE($5, doctor),
  updated_at = NOW()
WHERE id = $6 AND user_id = $7
RETURNING *;

-- name: Delete_Diagnostic_Schedule :one
DELETE FROM diagnostic_schedules
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: Get_Diagnsotic_Schedules_By_Centre :many
SELECT * FROM diagnostic_schedules
WHERE diagnostic_centre_id = $1
ORDER BY schedule_time DESC
LIMIT $2 OFFSET $3;

-- name: Get_Diagnsotic_Schedule_By_Centre :one
SELECT * FROM diagnostic_schedules
WHERE id = $1 AND diagnostic_centre_id = $2;


-- name: Update_Diagnostic_Schedule_By_Centre :one
UPDATE diagnostic_schedules
SET
  acceptance_status = COALESCE($1, acceptance_status),
  updated_at = NOW()
WHERE id = $2 AND diagnostic_centre_id = $3
RETURNING *;

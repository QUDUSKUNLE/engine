-- name: Create_Availability :many
WITH availability_params AS (
    SELECT 
        unnest($1::uuid[]) as diagnostic_centre_id,
        unnest($2::text[]) as day_of_week,
        unnest($3::time[]) as start_time,
        unnest($4::time[]) as end_time,
        unnest($5::int[]) as max_appointments,
        unnest($6::int[]) as slot_duration,
        unnest($7::int[]) as break_time
)
INSERT INTO diagnostic_centre_availability (
    diagnostic_centre_id,
    day_of_week,
    start_time,
    end_time,
    max_appointments,
    slot_duration,
    break_time
) 
SELECT * FROM availability_params
RETURNING *;

-- name: Get_Availability :many
SELECT * FROM diagnostic_centre_availability
WHERE diagnostic_centre_id = $1
AND ($2::text IS NULL OR day_of_week = $2)
ORDER BY
    CASE day_of_week
        WHEN 'monday' THEN 1
        WHEN 'tuesday' THEN 2
        WHEN 'wednesday' THEN 3
        WHEN 'thursday' THEN 4
        WHEN 'friday' THEN 5
        WHEN 'saturday' THEN 6
        WHEN 'sunday' THEN 7
    END;

-- name: Get_Diagnostic_Availability :many
SELECT * FROM diagnostic_centre_availability
WHERE diagnostic_centre_id = $1
ORDER BY created_at ASC;


-- name: Update_Availability :one
UPDATE diagnostic_centre_availability
SET
    start_time = COALESCE($3, start_time),
    end_time = COALESCE($4, end_time),
    max_appointments = COALESCE($5, max_appointments),
    slot_duration = COALESCE($6, slot_duration),
    break_time = COALESCE($7, break_time),
    updated_at = CURRENT_TIMESTAMP
WHERE diagnostic_centre_id = $1
AND day_of_week = $2
RETURNING *;

-- name: Update_Many_Availability :many
WITH update_params AS (
    SELECT 
        unnest($1::uuid[]) as diagnostic_centre_id,
        unnest($2::text[]) as day_of_week,
        unnest($3::time[]) as start_time,
        unnest($4::time[]) as end_time,
        unnest($5::int[]) as max_appointments,
        unnest($6::int[]) as slot_duration,
        unnest($7::int[]) as break_time
)
UPDATE diagnostic_centre_availability AS dca
SET
    start_time = COALESCE(up.start_time, dca.start_time),
    end_time = COALESCE(up.end_time, dca.end_time),
    max_appointments = COALESCE(up.max_appointments, dca.max_appointments),
    slot_duration = COALESCE(up.slot_duration, dca.slot_duration),
    break_time = COALESCE(up.break_time, dca.break_time),
    updated_at = CURRENT_TIMESTAMP
FROM update_params up
WHERE dca.diagnostic_centre_id = up.diagnostic_centre_id
AND dca.day_of_week = up.day_of_week
RETURNING dca.*;

-- name: Delete_Availability :exec
DELETE FROM diagnostic_centre_availability
WHERE diagnostic_centre_id = $1
AND day_of_week = $2;

-- name: Create_Single_Availability :one
INSERT INTO diagnostic_centre_availability (
    diagnostic_centre_id,
    day_of_week,
    start_time,
    end_time,
    max_appointments,
    slot_duration,
    break_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;


-- name: Create_Availability :one
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
) RETURNING *;

-- name: Get_Availability :many
SELECT * FROM diagnostic_centre_availability
WHERE diagnostic_centre_id = $1
AND ($2::weekday IS NULL OR day_of_week = $2)
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

-- name: Delete_Availability :exec
DELETE FROM diagnostic_centre_availability
WHERE diagnostic_centre_id = $1
AND day_of_week = $2;

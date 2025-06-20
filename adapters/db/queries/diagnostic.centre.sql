-- Inserts a new diagnostic record into the diagnostic_centres table.
-- name: Create_Diagnostic_Centre :one
INSERT INTO diagnostic_centres (
  diagnostic_centre_name,
  latitude,
  longitude,
  address,
  contact,
  doctors,
  available_tests,
  created_by,
  admin_id
) VALUES  (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- Retrieves a single diagnostic record by its ID.
-- name: Get_Diagnostic_Centre :one
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE dc.id = $1
GROUP BY dc.id;

-- Retrieves all diagnostic records with pagination.
-- name: Retrieve_Diagnostic_Centres :many
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
GROUP BY dc.id
ORDER BY dc.created_at DESC
LIMIT $1 OFFSET $2; 
 
-- Updates a diagnostic centre by the owner.
-- name: Update_Diagnostic_Centre_ByOwner :one
UPDATE diagnostic_centres
SET
  diagnostic_centre_name = COALESCE($3, diagnostic_centre_name),
  latitude = COALESCE($4, latitude),
  longitude = COALESCE($5, longitude),
  address = COALESCE($6, address),
  contact = COALESCE($7, contact),
  doctors = COALESCE($8, doctors),
  available_tests = COALESCE($9, available_tests),
  admin_id = COALESCE($10, admin_id),
  updated_at = NOW()
WHERE id = $1 AND created_by = $2
RETURNING *;

-- Deletes a diagnosticCentre only by the created_by.
-- name: Delete_Diagnostic_Centre_ByOwner :one
DELETE FROM diagnostic_centres WHERE id = $1 AND created_by = $2 RETURNING *;

-- Retrieves all diagnostic records for a specific owner.
-- name: List_Diagnostic_Centres_ByOwner :many
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE dc.created_by = $1
GROUP BY dc.id
ORDER BY dc.created_at DESC
LIMIT $2 OFFSET $3;

-- Searches diagnostic_centres by name with pagination.
-- name: Search_Diagnostic_Centres :many
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE
  (dc.diagnostic_centre_name ILIKE '%' || $1 || '%' OR $1 IS NULL)
  AND (dc.doctors @> $2 OR $2 IS NULL)
  AND (dc.available_tests @> $3 OR $3 IS NULL)
ORDER BY dc.created_at DESC
LIMIT $4 OFFSET $5;


-- GetADiagnosticCentreByOwner :one
-- name: Get_Diagnostic_Centre_ByOwner :one
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE dc.id = $1 AND dc.created_by = $2
GROUP BY dc.id;

-- GetDiagnosticCentreByManager
-- name: Get_Diagnostic_Centre_ByManager :one
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE dc.id = $1 AND dc.admin_id = $2
GROUP BY dc.id;

-- SearchDiagnosticWith Doctor type
-- name: Search_Diagnostic_Centres_ByDoctor :many
SELECT 
  dc.*,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE
  (dc.diagnostic_centre_name ILIKE '%' || $1 || '%' OR $1 IS NULL)
  AND (dc.doctors @> $2)
GROUP BY dc.id
ORDER BY dc.created_at DESC
LIMIT $3 OFFSET $4;

-- Retrieves the nearest diagnostic centres based on latitude and longitude.
-- name: Get_Nearest_Diagnostic_Centres :many
SELECT
  dc.id,
  dc.diagnostic_centre_name,
  dc.latitude,
  dc.longitude,
  dc.address,
  dc.contact,
  dc.doctors,
  dc.available_tests,
  dc.created_at,
  dc.updated_at,
  ARRAY_AGG(
    DISTINCT jsonb_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  CAST(
    6371 * acos(
      cos(radians($1)) * cos(radians(dc.latitude)) *
      cos(radians(dc.longitude) - radians($2)) +
      sin(radians($1)) * sin(radians(dc.latitude))
    ) AS DOUBLE PRECISION
  ) AS distance_km
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE
  dc.latitude IS NOT NULL
  AND dc.longitude IS NOT NULL
  AND (dc.doctors @> $3 OR $3 IS NULL)
  AND (dc.available_tests @> $4 OR $4 IS NULL)
GROUP BY
  dc.id
ORDER BY
  distance_km ASC
LIMIT 50;

-- name: Find_Nearest_Diagnostic_Centres_WhenRejected :many
SELECT
  dc.id,
  dc.diagnostic_centre_name,
  dc.latitude,
  dc.longitude,
  dc.address,
  dc.contact,
  dc.doctors,
  dc.available_tests,
  dc.created_at,
  dc.updated_at,
  ARRAY_AGG(
    json_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  CAST(
    6371 * acos(
      cos(radians($1)) * cos(radians(dc.latitude)) *
      cos(radians(dc.longitude) - radians($2)) +
      sin(radians($1)) * sin(radians(dc.latitude))
    ) AS DOUBLE PRECISION
  ) AS distance_km 
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
WHERE
  dc.id != $3 -- Exclude the current diagnostic centre
  AND dc.latitude IS NOT NULL
  AND dc.longitude IS NOT NULL
  AND (dc.doctors @> $4 OR $4 IS NULL) -- doctor type
  AND (dc.available_tests @> $5 OR $5 IS NULL) -- test type
GROUP BY
  dc.id
ORDER BY
  distance_km ASC
LIMIT 3;

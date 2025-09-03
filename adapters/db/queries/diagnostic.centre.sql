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

-- Retrieves a single diagnostic record by its ID and admin.
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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE dc.id = $1 AND dc.admin_id = $2
GROUP BY dc.id, prices.test_prices;

-- Retrieves a single diagnostic record by its ID and Owner.
-- name: Get_Diagnostic_Centre_By_Owner :one
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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE dc.id = $1 AND dc.created_by = $2
GROUP BY dc.id, prices.test_prices;

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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
GROUP BY dc.id, prices.test_prices
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
  admin_id = CASE 
    WHEN $10 IS NULL THEN NULL 
    ELSE COALESCE($10, admin_id)
  END,
  admin_assigned_by = CASE
    WHEN $10 IS NULL THEN NULL -- If admin_id is null, also null the assigned_by
    ELSE COALESCE($11, admin_assigned_by)
  END,
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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE dc.created_by = $1
GROUP BY dc.id, prices.test_prices
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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE
  (dc.diagnostic_centre_name ILIKE '%' || $1 || '%' OR $1 IS NULL)
  AND (dc.doctors @> $2 OR $2 IS NULL)
  AND (dc.available_tests @> $3 OR $3 IS NULL)
GROUP BY dc.id, prices.test_prices
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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE dc.id = $1 AND dc.created_by = $2
GROUP BY dc.id, prices.test_prices;

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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE dc.id = $1 AND dc.admin_id = $2
GROUP BY dc.id, prices.test_prices;

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
  ) FILTER (WHERE dca.diagnostic_centre_id IS NOT NULL) as availability,
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE
  (dc.diagnostic_centre_name ILIKE '%' || $1 || '%' OR $1 IS NULL)
  AND (dc.doctors @> $2)
GROUP BY dc.id, prices.test_prices
ORDER BY dc.created_at DESC
LIMIT $3 OFFSET $4;

-- Retrieves the nearest diagnostic centres based on latitude and longitude.
-- name: Get_Nearest_Diagnostic_Centres :many
WITH filtered_centres AS (
  SELECT dc.id
  FROM diagnostic_centres dc
  WHERE
    dc.latitude IS NOT NULL
    AND dc.longitude IS NOT NULL
    AND (dc.doctors @> $3 OR $3 IS NULL)
    AND (
      $4 = '' OR EXISTS (
        SELECT 1
        FROM diagnostic_centre_availability dca2
        WHERE dca2.diagnostic_centre_id = dc.id
        AND dca2.day_of_week = $4
      )
    )
    AND (
      $5 = '' OR EXISTS (
        SELECT 1
        FROM diagnostic_centre_test_prices dctp
        WHERE dctp.diagnostic_centre_id = dc.id
        AND dctp.test_type = $5
      )
    )
)
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

  -- Availability per day (if filtered or not)
  ARRAY_AGG(
    jsonb_build_object(
      'day_of_week', dca.day_of_week,
      'start_time', dca.start_time,
      'end_time', dca.end_time,
      'max_appointments', dca.max_appointments,
      'slot_duration', dca.slot_duration,
      'break_time', dca.break_time
    )
  ) FILTER (
    WHERE dca.day_of_week IS NOT NULL
    AND ($4 = '' OR dca.day_of_week = $4)
  ) AS availability,

  -- Strictly filtered test prices
  (
    SELECT jsonb_agg(
      jsonb_build_object(
        'test_type', dctp.test_type,
        'price', dctp.price
      )
    )
    FROM diagnostic_centre_test_prices dctp
    WHERE dctp.diagnostic_centre_id = dc.id
      AND ($5 = '' OR dctp.test_type = $5)
  ) AS test_prices,

  -- Distance calculation
  CAST(
    6371 * acos(
      cos(radians($1)) * cos(radians(dc.latitude)) *
      cos(radians(dc.longitude) - radians($2)) +
      sin(radians($1)) * sin(radians(dc.latitude))
    ) AS DOUBLE PRECISION
  ) AS distance_km

FROM filtered_centres fc
JOIN diagnostic_centres dc ON dc.id = fc.id
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id

GROUP BY
  dc.id, dc.diagnostic_centre_name, dc.latitude, dc.longitude,
  dc.address, dc.contact, dc.doctors, dc.available_tests,
  dc.created_at, dc.updated_at

ORDER BY distance_km ASC
LIMIT $6 OFFSET $7;


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
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices,
  CAST(
    6371 * acos(
      cos(radians($1)) * cos(radians(dc.latitude)) *
      cos(radians(dc.longitude) - radians($2)) +
      sin(radians($1)) * sin(radians(dc.latitude))
    ) AS DOUBLE PRECISION
  ) AS distance_km 
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE
  dc.id != $3 -- Exclude the current diagnostic centre
  AND dc.latitude IS NOT NULL
  AND dc.longitude IS NOT NULL
  AND (dc.doctors @> $4 OR $4 IS NULL) -- doctor type
  AND (dc.available_tests @> $5 OR $5 IS NULL) -- test type
GROUP BY
  dc.id, prices.test_prices
ORDER BY
  distance_km ASC
LIMIT 3;


-- name: Get_Diagnostic_Centre_Managers :many
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
  COALESCE(prices.test_prices, '[]'::jsonb) AS test_prices
FROM diagnostic_centres dc
LEFT JOIN diagnostic_centre_availability dca ON dc.id = dca.diagnostic_centre_id
LEFT JOIN LATERAL (
  SELECT jsonb_agg(
    jsonb_build_object(
      'test_type', dctp.test_type,
      'price', dctp.price
    )
  ) AS test_prices
  FROM diagnostic_centre_test_prices dctp
  WHERE dctp.diagnostic_centre_id = dc.id
) prices ON true
WHERE
  dc.admin_id = $1
GROUP BY
  dc.id, prices.test_prices
ORDER BY dc.created_at DESC
LIMIT $2 OFFSET $3;


-- name: UnassignAdmin :one
UPDATE diagnostic_centres
SET 
  admin_id = NULL,
  admin_unassigned_by = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AssignAdmin :one
UPDATE diagnostic_centres
SET 
  admin_id = $2,
  admin_assigned_by = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetAdminHistory :many
SELECT 
  dc.id,
  dc.diagnostic_centre_name,
  dc.admin_id,
  u.fullname as admin_name,
  dc.admin_assigned_at,
  dc.admin_unassigned_at,
  dc.admin_status
FROM diagnostic_centres dc
LEFT JOIN users u ON u.id = dc.admin_id
WHERE dc.id = $1
ORDER BY dc.admin_assigned_at DESC;

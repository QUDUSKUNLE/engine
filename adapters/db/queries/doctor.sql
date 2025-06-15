-- name: GetSampleDoctor :one
-- This is just a helper query to make sqlc generate the Doctor type
SELECT 'Male'::doctor as doctor;

-- name: GetSampleDoctorArray :one
-- This is just a helper query to make sqlc generate the Doctor type as array
SELECT ARRAY['Male']::doctor[] as doctors;

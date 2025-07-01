-- name: Create_Test_Price :many
WITH test_price_params AS (
    SELECT 
        unnest($1::uuid[]) as diagnostic_centre_id,
        unnest($2::text[]) as test_type,
        unnest($3::float[]) as price,
        unnest($4::varchar[]) as currency,
        unnest($5::boolean[]) as is_active
)
INSERT INTO diagnostic_centre_test_prices (
    diagnostic_centre_id,
    test_type,
    price,
    currency,
    is_active
) 
SELECT * FROM test_price_params
RETURNING *;

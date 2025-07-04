-- name: Create_Payment :one
INSERT INTO payments (
    appointment_id,
    patient_id,
    diagnostic_centre_id,
    amount,
    currency,
    payment_method,
    payment_metadata,
    payment_provider,
    provider_reference,
    provider_metadata,
    transaction_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: Get_Payment :one
SELECT * FROM payments WHERE id = $1;

-- name: List_Payments :many
SELECT * FROM payments 
WHERE 
    (CASE WHEN $1::UUID IS NOT NULL THEN diagnostic_centre_id = $1 ELSE TRUE END) AND
    (CASE WHEN $2::UUID IS NOT NULL THEN patient_id = $2 ELSE TRUE END) AND
    (CASE WHEN $3::payment_status IS NOT NULL THEN payment_status = $3 ELSE TRUE END) AND
    (CASE WHEN $4::TIMESTAMP IS NOT NULL THEN payment_date >= $4 ELSE TRUE END) AND
    (CASE WHEN $5::TIMESTAMP IS NOT NULL THEN payment_date <= $5 ELSE TRUE END)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: Update_Payment_Status :one
UPDATE payments 
SET 
    payment_status = COALESCE($2, payment_status),
    transaction_id = COALESCE($3, transaction_id),
    payment_metadata = COALESCE($4, payment_metadata),
    payment_date = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
RETURNING *;

-- name: Refund_Payment :one
UPDATE payments 
SET 
    payment_status = 'refunded',
    refund_amount = $2,
    refund_reason = $3,
    refund_date = CURRENT_TIMESTAMP,
    refunded_by = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    payments.id = $1 AND 
    payment_status = 'success' AND
    NOT EXISTS (
        SELECT 1 FROM payments AS p
        WHERE p.id = payments.id AND p.refund_amount IS NOT NULL
    )
RETURNING *;

-- name: GetPaymentByReference :one
SELECT * FROM payments WHERE provider_reference = $1 LIMIT 1;

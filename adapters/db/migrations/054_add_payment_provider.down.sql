-- Drop added columns
ALTER TABLE payments 
    DROP COLUMN provider_metadata,
    DROP COLUMN provider_reference,
    DROP COLUMN payment_provider;

-- Drop the enum type
DROP TYPE payment_provider;

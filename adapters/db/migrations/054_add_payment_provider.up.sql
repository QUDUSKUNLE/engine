-- Create payment provider enum
CREATE TYPE payment_provider AS ENUM (
    'PAYSTACK',
    'FLUTTERWAVE',
    'STRIPE',
    'MONNIFY'
);

-- Add payment_provider column to payments table
ALTER TABLE payments 
    ADD COLUMN payment_provider payment_provider NOT NULL DEFAULT 'PAYSTACK',
    ADD COLUMN provider_reference VARCHAR(255),
    ADD COLUMN provider_metadata JSONB;

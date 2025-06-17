-- Create payment method ENUM
CREATE TYPE payment_method AS ENUM (
    'card',
    'transfer',
    'cash',
    'wallet'
);

-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Foreign keys
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    patient_id UUID NOT NULL REFERENCES users(id),
    diagnostic_centre_id UUID NOT NULL REFERENCES diagnostic_centres(id),
    
    -- Payment details
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    payment_method payment_method NOT NULL,
    payment_status payment_status NOT NULL DEFAULT 'pending',
    
    -- Transaction details
    transaction_id TEXT NULL,
    payment_metadata JSONB NULL,
    payment_date TIMESTAMP WITH TIME ZONE NULL,
    
    -- Refund info
    refund_amount DECIMAL(10,2) NULL CHECK (refund_amount IS NULL OR (refund_amount > 0 AND refund_amount <= amount)),
    refund_reason TEXT NULL,
    refund_date TIMESTAMP WITH TIME ZONE NULL,
    refunded_by UUID NULL REFERENCES users(id),
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CHECK ((payment_status = 'refunded' AND refund_amount IS NOT NULL AND refund_reason IS NOT NULL AND refund_date IS NOT NULL AND refunded_by IS NOT NULL) OR
           (payment_status != 'refunded' AND refund_amount IS NULL AND refund_reason IS NULL AND refund_date IS NULL AND refunded_by IS NULL)),
    CHECK ((payment_status = 'success' AND payment_date IS NOT NULL AND transaction_id IS NOT NULL) OR
           (payment_status != 'success' AND (payment_date IS NULL OR transaction_id IS NULL)))
);

-- Add indexes for common queries
CREATE INDEX idx_payments_appointment ON payments(appointment_id);
CREATE INDEX idx_payments_patient ON payments(patient_id);
CREATE INDEX idx_payments_centre ON payments(diagnostic_centre_id);
CREATE INDEX idx_payments_status ON payments(payment_status);
CREATE INDEX idx_payments_date ON payments(payment_date DESC NULLS LAST);
CREATE INDEX idx_payments_created ON payments(created_at DESC);

-- Function to update updated_at timestamp
CREATE TRIGGER update_payment_timestamp BEFORE UPDATE
ON payments FOR EACH ROW EXECUTE FUNCTION
update_updated_at_column();

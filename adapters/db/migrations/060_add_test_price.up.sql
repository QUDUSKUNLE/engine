CREATE TABLE IF NOT EXISTS diagnostic_centre_test_prices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  diagnostic_centre_id UUID NOT NULL REFERENCES diagnostic_centres(id) ON DELETE CASCADE,
  test_type TEXT NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  UNIQUE (diagnostic_centre_id, test_type)
);

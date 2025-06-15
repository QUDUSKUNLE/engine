-- Enable the btree_gist extension which is needed for time range operations
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Create function for automatically updating updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

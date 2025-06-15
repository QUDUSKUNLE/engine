-- Drop the function first
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop the btree_gist extension
DROP EXTENSION IF EXISTS btree_gist;

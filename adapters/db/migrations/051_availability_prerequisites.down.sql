-- Drop prerequisites in reverse order
DROP TYPE IF EXISTS timerange;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
DROP EXTENSION IF EXISTS btree_gist;

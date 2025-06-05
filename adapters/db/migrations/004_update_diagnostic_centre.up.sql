ALTER TYPE user_enum RENAME VALUE 'DIAGNOSTIC_CENTRE' TO 'DIAGNOSTIC_CENTRE_OWNER';

ALTER TYPE user_enum RENAME VALUE 'DIAGNOSTIC_MANAGER' TO 'DIAGNOSTIC_CENTRE_MANAGER';

-- CREATE TABLE IF NOT EXISTS diagnostic_manager (
--   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   user_id UUID NOT NULL REFERENCES users(id),
--   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--   UNIQUE (user_id, diagnostic_centre_name),
-- );

-- CREATE INDEX idx_diagnostics_user_id ON diagnostics(user_id);

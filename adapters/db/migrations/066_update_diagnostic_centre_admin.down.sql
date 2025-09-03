-- First drop the trigger
DROP TRIGGER IF EXISTS track_admin_changes ON diagnostic_centres;

-- Then drop the trigger function
DROP FUNCTION IF EXISTS track_admin_changes();

-- Finally remove the added columns and restore NOT NULL constraint
ALTER TABLE diagnostic_centres
    DROP COLUMN IF EXISTS admin_status,
    DROP COLUMN IF EXISTS admin_unassigned_by,
    DROP COLUMN IF EXISTS admin_unassigned_at,
    DROP COLUMN IF EXISTS admin_assigned_by,
    DROP COLUMN IF EXISTS admin_assigned_at,
    ALTER COLUMN admin_id SET NOT NULL;

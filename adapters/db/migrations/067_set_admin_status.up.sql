-- Update existing records to set proper admin status and timestamps
UPDATE diagnostic_centres
SET 
    admin_status = CASE 
        WHEN admin_id IS NULL THEN 'UNASSIGNED'
        ELSE 'ASSIGNED'
    END,
    admin_assigned_at = CASE 
        WHEN admin_id IS NOT NULL THEN created_at
        ELSE NULL
    END;


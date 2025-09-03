-- Update existing records to set admin_assigned_by to created_by
UPDATE diagnostic_centres
SET 
    admin_assigned_by = created_by,
    admin_assigned_at = created_at
WHERE 
    admin_id IS NOT NULL 
    AND admin_assigned_by IS NULL;

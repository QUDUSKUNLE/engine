-- Reset admin_assigned_by to NULL for records where it matches created_by
UPDATE diagnostic_centres
SET 
    admin_assigned_by = NULL,
    admin_assigned_at = NULL
WHERE 
    admin_assigned_by = created_by;

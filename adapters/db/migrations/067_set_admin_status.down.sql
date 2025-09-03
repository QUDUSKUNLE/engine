-- Reset the admin status fields
UPDATE diagnostic_centres
SET 
    admin_status = 'ASSIGNED',
    admin_assigned_at = NULL,
    admin_unassigned_at = NULL;

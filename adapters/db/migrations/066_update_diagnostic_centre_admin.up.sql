-- Add admin tracking columns
ALTER TABLE diagnostic_centres
    ALTER COLUMN admin_id DROP NOT NULL,
    ADD COLUMN admin_assigned_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN admin_assigned_by UUID REFERENCES users(id),
    ADD COLUMN admin_unassigned_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN admin_unassigned_by UUID REFERENCES users(id),
    ADD COLUMN admin_status VARCHAR(20) DEFAULT 'ASSIGNED' CHECK (admin_status IN ('ASSIGNED', 'UNASSIGNED'));


CREATE OR REPLACE FUNCTION track_admin_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.admin_id IS NOT NULL AND NEW.admin_id IS NULL THEN
        -- Admin being unassigned
        NEW.admin_unassigned_at = NOW();
        NEW.admin_status = 'UNASSIGNED';
    ELSIF NEW.admin_id IS NOT NULL AND (OLD.admin_id IS NULL OR OLD.admin_id != NEW.admin_id) THEN
        -- New admin being assigned
        NEW.admin_assigned_at = NOW();
        NEW.admin_status = 'ASSIGNED';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER track_admin_changes
    BEFORE UPDATE ON diagnostic_centres
    FOR EACH ROW
    EXECUTE FUNCTION track_admin_changes();

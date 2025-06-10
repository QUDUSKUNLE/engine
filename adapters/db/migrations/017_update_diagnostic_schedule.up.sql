ALTER TABLE diagnostic_schedules
    ADD CONSTRAINT diagnostic_schedules_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);

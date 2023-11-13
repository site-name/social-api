ALTER TABLE ONLY export_events
    ADD CONSTRAINT fk_export_events_export_files FOREIGN KEY (export_file_id) REFERENCES export_files(id);
ALTER TABLE ONLY export_events
    ADD CONSTRAINT fk_export_events_users FOREIGN KEY (user_id) REFERENCES users(id);
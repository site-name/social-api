ALTER TABLE ONLY export_events
    ADD CONSTRAINT fk_export_events_export_files FOREIGN KEY (exportfileid) REFERENCES export_files(id);
ALTER TABLE ONLY export_events
    ADD CONSTRAINT fk_export_events_users FOREIGN KEY (userid) REFERENCES users(id);

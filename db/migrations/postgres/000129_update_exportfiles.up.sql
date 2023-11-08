ALTER TABLE ONLY export_files
    ADD CONSTRAINT fk_export_files_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

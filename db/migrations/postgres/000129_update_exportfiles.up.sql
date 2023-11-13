ALTER TABLE ONLY export_files
    ADD CONSTRAINT fk_export_files_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
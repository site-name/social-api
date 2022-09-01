ALTER TABLE ONLY exportfiles
    ADD CONSTRAINT fk_exportfiles_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

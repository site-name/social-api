ALTER TABLE ONLY staffnotificationrecipients
    ADD CONSTRAINT fk_staffnotificationrecipients_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

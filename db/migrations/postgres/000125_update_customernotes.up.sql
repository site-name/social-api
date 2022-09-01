ALTER TABLE ONLY customernotes
    ADD CONSTRAINT fk_customernotes_users FOREIGN KEY (userid) REFERENCES users(id);

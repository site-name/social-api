ALTER TABLE ONLY customer_notes
    ADD CONSTRAINT fk_customer_notes_users FOREIGN KEY (user_id) REFERENCES users(id);
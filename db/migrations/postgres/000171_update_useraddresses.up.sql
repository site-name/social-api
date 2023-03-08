ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT fk_useraddresses_addresses FOREIGN KEY (addressid) REFERENCES addresses(id) ON DELETE CASCADE;
ALTER TABLE ONLY useraddresses
    ADD CONSTRAINT fk_useraddresses_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

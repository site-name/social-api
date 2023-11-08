ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT fk_user_addresses_addresses FOREIGN KEY (addressid) REFERENCES addresses(id) ON DELETE CASCADE;
ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT fk_user_addresses_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

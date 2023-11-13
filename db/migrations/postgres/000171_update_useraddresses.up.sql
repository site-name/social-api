ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT fk_user_addresses_addresses FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE;
ALTER TABLE ONLY user_addresses
    ADD CONSTRAINT fk_user_addresses_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

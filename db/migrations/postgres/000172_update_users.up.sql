ALTER TABLE ONLY users
    ADD CONSTRAINT fk_users_addresses FOREIGN KEY (default_shipping_address_id) REFERENCES addresses(id);

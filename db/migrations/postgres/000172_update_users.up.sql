ALTER TABLE ONLY users
    ADD CONSTRAINT fk_users_shipping_addresses FOREIGN KEY (default_shipping_address_id) REFERENCES addresses(id);

ALTER TABLE ONLY users
    ADD CONSTRAINT fk_users_billing_addresses FOREIGN KEY (default_billing_address_id) REFERENCES addresses(id);


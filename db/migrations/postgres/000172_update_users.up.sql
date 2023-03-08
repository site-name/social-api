ALTER TABLE ONLY users
    ADD CONSTRAINT fk_users_addresses FOREIGN KEY (defaultshippingaddressid) REFERENCES addresses(id);

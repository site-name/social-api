ALTER TABLE ONLY warehouses
    ADD CONSTRAINT fk_warehouses_addresses FOREIGN KEY (address_id) REFERENCES addresses(id);

ALTER TABLE ONLY warehouses
    ADD CONSTRAINT fk_warehouses_addresses FOREIGN KEY (addressid) REFERENCES addresses(id);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id) ON DELETE CASCADE;

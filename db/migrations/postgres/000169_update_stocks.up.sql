ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_product_variants FOREIGN KEY (productvariantid) REFERENCES product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_warehouses FOREIGN KEY (warehouseid) REFERENCES warehouses(id) ON DELETE CASCADE;

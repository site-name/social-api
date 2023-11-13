ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_product_variants FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY stocks
    ADD CONSTRAINT fk_stocks_warehouses FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE;

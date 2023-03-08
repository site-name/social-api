ALTER TABLE ONLY productvariants
    ADD CONSTRAINT fk_productvariants_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

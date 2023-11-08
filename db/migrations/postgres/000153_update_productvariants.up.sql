ALTER TABLE ONLY product_variants
    ADD CONSTRAINT fk_product_variants_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE ONLY product_translations
    ADD CONSTRAINT fk_product_translations_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

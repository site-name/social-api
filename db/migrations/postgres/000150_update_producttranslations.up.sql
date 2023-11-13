ALTER TABLE ONLY product_translations
    ADD CONSTRAINT fk_product_translations_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
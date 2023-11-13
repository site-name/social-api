ALTER TABLE ONLY product_media
    ADD CONSTRAINT fk_product_media_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
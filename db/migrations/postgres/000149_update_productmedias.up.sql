ALTER TABLE ONLY product_media
    ADD CONSTRAINT fk_product_media_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

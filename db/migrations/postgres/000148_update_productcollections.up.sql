ALTER TABLE ONLY product_collections
    ADD CONSTRAINT fk_product_collections_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY product_collections
    ADD CONSTRAINT fk_product_collections_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

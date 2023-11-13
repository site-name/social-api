ALTER TABLE ONLY product_collections
    ADD CONSTRAINT fk_product_collections_collections FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY product_collections
    ADD CONSTRAINT fk_product_collections_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
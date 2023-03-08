ALTER TABLE ONLY productcollections
    ADD CONSTRAINT fk_productcollections_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;
ALTER TABLE ONLY productcollections
    ADD CONSTRAINT fk_productcollections_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

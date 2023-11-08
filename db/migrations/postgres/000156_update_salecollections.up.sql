ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT fk_sale_collections_collections FOREIGN KEY (collectionid) REFERENCES collections(id);
ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT fk_sale_collections_sales FOREIGN KEY (saleid) REFERENCES sales(id);

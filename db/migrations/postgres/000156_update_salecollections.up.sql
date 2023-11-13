ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT fk_sale_collections_collections FOREIGN KEY (collection_id) REFERENCES collections(id);
ALTER TABLE ONLY sale_collections
    ADD CONSTRAINT fk_sale_collections_sales FOREIGN KEY (sale_id) REFERENCES sales(id);

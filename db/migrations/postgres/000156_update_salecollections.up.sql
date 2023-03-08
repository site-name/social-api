ALTER TABLE ONLY salecollections
    ADD CONSTRAINT fk_salecollections_collections FOREIGN KEY (collectionid) REFERENCES collections(id);
ALTER TABLE ONLY salecollections
    ADD CONSTRAINT fk_salecollections_sales FOREIGN KEY (saleid) REFERENCES sales(id);

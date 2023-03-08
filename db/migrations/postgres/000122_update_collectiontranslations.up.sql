ALTER TABLE ONLY collectiontranslations
    ADD CONSTRAINT fk_collectiontranslations_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;

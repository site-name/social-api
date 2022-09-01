ALTER TABLE ONLY collections
    ADD CONSTRAINT fk_collections_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;

ALTER TABLE ONLY collection_translations
    ADD CONSTRAINT fk_collection_translations_collections FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE;
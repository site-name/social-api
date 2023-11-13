ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT fk_collection_channel_listings_channels FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT fk_collection_channel_listings_collections FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE;
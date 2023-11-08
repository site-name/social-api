ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT fk_collection_channel_listings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY collection_channel_listings
    ADD CONSTRAINT fk_collection_channel_listings_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;

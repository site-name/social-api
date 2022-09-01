ALTER TABLE ONLY collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY collectionchannellistings
    ADD CONSTRAINT fk_collectionchannellistings_collections FOREIGN KEY (collectionid) REFERENCES collections(id) ON DELETE CASCADE;

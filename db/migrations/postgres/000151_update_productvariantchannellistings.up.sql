ALTER TABLE ONLY productvariantchannellistings
    ADD CONSTRAINT fk_productvariantchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY productvariantchannellistings
    ADD CONSTRAINT fk_productvariantchannellistings_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;

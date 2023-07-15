ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id) ON DELETE CASCADE;

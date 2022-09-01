ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_channels FOREIGN KEY (channelid) REFERENCES channels(id);
ALTER TABLE ONLY shippingzonechannels
    ADD CONSTRAINT fk_shippingzonechannels_shippingzones FOREIGN KEY (shippingzoneid) REFERENCES shippingzones(id);

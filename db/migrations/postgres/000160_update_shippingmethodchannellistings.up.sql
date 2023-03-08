ALTER TABLE ONLY shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY shippingmethodchannellistings
    ADD CONSTRAINT fk_shippingmethodchannellistings_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id) ON DELETE CASCADE;

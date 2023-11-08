ALTER TABLE ONLY shipping_zone_channels
    ADD CONSTRAINT fk_shipping_zone_channels_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY shipping_zone_channels
    ADD CONSTRAINT fk_shipping_zone_channels_shipping_zones FOREIGN KEY (shippingzoneid) REFERENCES shipping_zones(id) ON DELETE CASCADE;

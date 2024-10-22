ALTER TABLE ONLY shipping_method_channel_listings
    ADD CONSTRAINT fk_shipping_method_channel_listings_channels FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY shipping_method_channel_listings
    ADD CONSTRAINT fk_shipping_method_channel_listings_shipping_methods FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id) ON DELETE CASCADE;
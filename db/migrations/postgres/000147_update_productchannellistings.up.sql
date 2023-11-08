ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT fk_product_channel_listings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY product_channel_listings
    ADD CONSTRAINT fk_product_channel_listings_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT fk_sale_channel_listings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT fk_sale_channel_listings_sales FOREIGN KEY (saleid) REFERENCES sales(id) ON DELETE CASCADE;

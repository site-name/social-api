ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT fk_sale_channel_listings_channels FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY sale_channel_listings
    ADD CONSTRAINT fk_sale_channel_listings_sales FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE;

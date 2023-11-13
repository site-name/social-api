ALTER TABLE ONLY voucher_channel_listings
    ADD CONSTRAINT fk_voucher_channel_listings_channels FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucher_channel_listings
    ADD CONSTRAINT fk_voucher_channel_listings_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;

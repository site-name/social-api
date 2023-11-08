ALTER TABLE ONLY product_variant_channel_listings
    ADD CONSTRAINT fk_product_variant_channel_listings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY product_variant_channel_listings
    ADD CONSTRAINT fk_product_variant_channel_listings_product_variants FOREIGN KEY (variantid) REFERENCES product_variants(id) ON DELETE CASCADE;

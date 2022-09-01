ALTER TABLE ONLY productchannellistings
    ADD CONSTRAINT fk_productchannellistings_channels FOREIGN KEY (channelid) REFERENCES channels(id) ON DELETE CASCADE;
ALTER TABLE ONLY productchannellistings
    ADD CONSTRAINT fk_productchannellistings_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

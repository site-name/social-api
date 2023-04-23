ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_addresses FOREIGN KEY (billingaddressid) REFERENCES addresses(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_channels FOREIGN KEY (channelid) REFERENCES channels(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_warehouses FOREIGN KEY (collectionpointid) REFERENCES warehouses(id);

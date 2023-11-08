ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_addresses FOREIGN KEY (billingaddressid) REFERENCES addresses(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_channels FOREIGN KEY (channelid) REFERENCES channels(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_shipping_methods FOREIGN KEY (shippingmethodid) REFERENCES shipping_methods(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_warehouses FOREIGN KEY (collectionpointid) REFERENCES warehouses(id);

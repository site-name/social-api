ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_addresses FOREIGN KEY (billingaddressid) REFERENCES addresses(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_channels FOREIGN KEY (channelid) REFERENCES channels(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_orders FOREIGN KEY (originalid) REFERENCES orders(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_shippingmethods FOREIGN KEY (shippingmethodid) REFERENCES shippingmethods(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_users FOREIGN KEY (userid) REFERENCES users(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_warehouses FOREIGN KEY (collectionpointid) REFERENCES warehouses(id);

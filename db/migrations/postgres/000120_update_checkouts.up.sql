ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_addresses FOREIGN KEY (billing_address_id) REFERENCES addresses(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_channels FOREIGN KEY (channel_id) REFERENCES channels(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_shipping_methods FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id);
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE ONLY checkouts
    ADD CONSTRAINT fk_checkouts_warehouses FOREIGN KEY (collection_point_id) REFERENCES warehouses(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_addresses FOREIGN KEY (billing_address_id) REFERENCES addresses(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_channels FOREIGN KEY (channel_id) REFERENCES channels(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_orders FOREIGN KEY (original_id) REFERENCES orders(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_shipping_methods FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_users FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_vouchers FOREIGN KEY (voucher_id) REFERENCES vouchers(id);
ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_warehouses FOREIGN KEY (collection_point_id) REFERENCES warehouses(id);
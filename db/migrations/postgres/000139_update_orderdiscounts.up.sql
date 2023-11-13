ALTER TABLE ONLY order_discounts
    ADD CONSTRAINT fk_order_discounts_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
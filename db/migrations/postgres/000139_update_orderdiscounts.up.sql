ALTER TABLE ONLY order_discounts
    ADD CONSTRAINT fk_order_discounts_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;

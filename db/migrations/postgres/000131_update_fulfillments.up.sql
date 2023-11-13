ALTER TABLE ONLY fulfillments
    ADD CONSTRAINT fk_fulfillments_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
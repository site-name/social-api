ALTER TABLE ONLY fulfillments
    ADD CONSTRAINT fk_fulfillments_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE ONLY order_events
    ADD CONSTRAINT fk_order_events_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY order_events
    ADD CONSTRAINT fk_order_events_users FOREIGN KEY (userid) REFERENCES users(id);

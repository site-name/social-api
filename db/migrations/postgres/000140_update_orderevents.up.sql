ALTER TABLE ONLY order_events
    ADD CONSTRAINT fk_order_events_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY order_events
    ADD CONSTRAINT fk_order_events_users FOREIGN KEY (user_id) REFERENCES users(id);
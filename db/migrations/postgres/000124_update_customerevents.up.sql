ALTER TABLE ONLY customer_events
    ADD CONSTRAINT fk_customer_events_orders FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE ONLY customer_events
    ADD CONSTRAINT fk_customer_events_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE ONLY customer_events
    ADD CONSTRAINT fk_customer_events_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY customer_events
    ADD CONSTRAINT fk_customer_events_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY orderevents
    ADD CONSTRAINT fk_orderevents_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;
ALTER TABLE ONLY orderevents
    ADD CONSTRAINT fk_orderevents_users FOREIGN KEY (userid) REFERENCES users(id);

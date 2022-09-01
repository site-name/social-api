ALTER TABLE ONLY customerevents
    ADD CONSTRAINT fk_customerevents_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY customerevents
    ADD CONSTRAINT fk_customerevents_users FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;

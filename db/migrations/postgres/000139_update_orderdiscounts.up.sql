ALTER TABLE ONLY orderdiscounts
    ADD CONSTRAINT fk_orderdiscounts_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;

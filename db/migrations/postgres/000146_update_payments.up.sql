ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token);
ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_orders FOREIGN KEY (orderid) REFERENCES orders(id);

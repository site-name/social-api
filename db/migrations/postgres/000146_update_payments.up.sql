ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_checkouts FOREIGN KEY (checkout_id) REFERENCES checkouts(token);
ALTER TABLE ONLY payments
    ADD CONSTRAINT fk_payments_orders FOREIGN KEY (order_id) REFERENCES orders(id);
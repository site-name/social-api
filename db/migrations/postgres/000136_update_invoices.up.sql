ALTER TABLE ONLY invoices
    ADD CONSTRAINT fk_invoices_orders FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE ONLY invoices
    ADD CONSTRAINT fk_invoices_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_invoices FOREIGN KEY (invoiceid) REFERENCES invoices(id);
ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_users FOREIGN KEY (userid) REFERENCES users(id);

ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_invoices FOREIGN KEY (invoice_id) REFERENCES invoices(id);
ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_orders FOREIGN KEY (order_id) REFERENCES orders(id);
ALTER TABLE ONLY invoice_events
    ADD CONSTRAINT fk_invoice_events_users FOREIGN KEY (user_id) REFERENCES users(id);
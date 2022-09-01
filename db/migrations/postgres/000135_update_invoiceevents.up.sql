ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_invoices FOREIGN KEY (invoiceid) REFERENCES invoices(id);
ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_orders FOREIGN KEY (orderid) REFERENCES orders(id);
ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_users FOREIGN KEY (userid) REFERENCES users(id);

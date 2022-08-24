CREATE TABLE IF NOT EXISTS invoiceevents (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  type character varying(255),
  invoiceid character varying(36),
  orderid character varying(36),
  userid character varying(36),
  parameters text
);

ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_invoices FOREIGN KEY (invoiceid) REFERENCES invoices(id);

ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_orders FOREIGN KEY (orderid) REFERENCES orders(id);

ALTER TABLE ONLY invoiceevents
    ADD CONSTRAINT fk_invoiceevents_users FOREIGN KEY (userid) REFERENCES users(id);

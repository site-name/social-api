CREATE TABLE IF NOT EXISTS invoices (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderid character varying(36),
  number character varying(255),
  createat bigint,
  externalurl character varying(2048),
  invoicefile text,
  metadata jsonb,
  privatemetadata jsonb
);
ALTER TABLE ONLY invoices
    ADD CONSTRAINT fk_invoices_orders FOREIGN KEY (orderid) REFERENCES orders(id);

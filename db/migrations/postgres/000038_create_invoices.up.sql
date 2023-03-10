CREATE TABLE IF NOT EXISTS invoices (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderid character varying(36),
  number character varying(255),
  createat bigint,
  externalurl character varying(2048),
  status character varying(50),
  message character varying(255),
  updateat bigint,
  invoicefile text,
  metadata jsonb,
  privatemetadata jsonb
);
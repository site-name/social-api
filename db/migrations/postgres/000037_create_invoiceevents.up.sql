CREATE TABLE IF NOT EXISTS invoice_events (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  type character varying(255),
  invoiceid character varying(36),
  orderid character varying(36),
  userid character varying(36),
  parameters text
);

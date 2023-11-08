CREATE TABLE IF NOT EXISTS order_events (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  type character varying(255),
  orderid character varying(36),
  parameters text,
  userid character varying(36)
);


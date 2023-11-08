CREATE TABLE IF NOT EXISTS customer_events (
  id character varying(36) NOT NULL PRIMARY KEY,
  date bigint,
  type character varying(255),
  orderid character varying(36),
  userid character varying(36),
  parameters jsonb
);


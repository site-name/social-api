CREATE TABLE IF NOT EXISTS fulfillmentlines (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderlineid character varying(36),
  fulfillmentid character varying(36),
  quantity integer,
  stockid character varying(36)
);

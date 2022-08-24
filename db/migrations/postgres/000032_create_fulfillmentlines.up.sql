CREATE TABLE IF NOT EXISTS fulfillmentlines (
  id character varying(36) NOT NULL PRIMARY KEY,
  orderlineid character varying(36),
  fulfillmentid character varying(36),
  quantity integer,
  stockid character varying(36)
);

ALTER TABLE ONLY fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_orderlines FOREIGN KEY (orderlineid) REFERENCES orderlines(id) ON DELETE CASCADE;

ALTER TABLE ONLY fulfillmentlines
    ADD CONSTRAINT fk_fulfillmentlines_stocks FOREIGN KEY (stockid) REFERENCES stocks(id);

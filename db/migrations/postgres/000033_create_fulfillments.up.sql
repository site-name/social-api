CREATE TABLE IF NOT EXISTS fulfillments (
  id character varying(36) NOT NULL PRIMARY KEY,
  fulfillmentorder integer,
  orderid character varying(36),
  status character varying(32),
  trackingnumber character varying(255),
  createat bigint,
  shippingrefundamount double precision,
  totalrefundamount double precision,
  metadata text,
  privatemetadata text
);

CREATE INDEX idx_fulfillments_status ON fulfillments USING btree (status);

CREATE INDEX idx_fulfillments_tracking_number ON fulfillments USING btree (trackingnumber);

ALTER TABLE ONLY fulfillments
    ADD CONSTRAINT fk_fulfillments_orders FOREIGN KEY (orderid) REFERENCES orders(id) ON DELETE CASCADE;


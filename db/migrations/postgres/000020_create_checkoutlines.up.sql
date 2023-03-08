CREATE TABLE IF NOT EXISTS checkoutlines (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  checkoutid character varying(36),
  variantid character varying(36),
  quantity integer
);

CREATE INDEX idx_checkoutlines_checkout_id ON checkoutlines USING btree (checkoutid);

CREATE INDEX idx_checkoutlines_variant_id ON checkoutlines USING btree (variantid);

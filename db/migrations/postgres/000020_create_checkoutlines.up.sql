CREATE TABLE IF NOT EXISTS checkoutlines (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  checkoutid character varying(36),
  variantid character varying(36),
  quantity integer
);

CREATE INDEX idx_checkoutlines_checkout_id ON checkoutlines USING btree (checkoutid);

CREATE INDEX idx_checkoutlines_variant_id ON checkoutlines USING btree (variantid);

ALTER TABLE ONLY checkoutlines
    ADD CONSTRAINT fk_checkoutlines_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token) ON DELETE CASCADE;

ALTER TABLE ONLY checkoutlines
    ADD CONSTRAINT fk_checkoutlines_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;

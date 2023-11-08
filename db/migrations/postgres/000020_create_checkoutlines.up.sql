CREATE TABLE IF NOT EXISTS checkout_lines (
  id character varying(36) NOT NULL PRIMARY KEY,
  createat bigint,
  checkoutid character varying(36),
  variantid character varying(36),
  quantity integer
);

CREATE INDEX idx_checkout_lines_checkout_id ON checkout_lines USING btree (checkoutid);

CREATE INDEX idx_checkout_lines_variant_id ON checkout_lines USING btree (variantid);

CREATE TABLE IF NOT EXISTS checkout_lines (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  checkout_id varchar(36) NOT NULL,
  variant_id varchar(36) NOT NULL,
  quantity integer NOT NULL
);

CREATE INDEX idx_checkout_lines_checkout_id ON checkout_lines USING btree (checkout_id);

CREATE INDEX idx_checkout_lines_variant_id ON checkout_lines USING btree (variant_id);
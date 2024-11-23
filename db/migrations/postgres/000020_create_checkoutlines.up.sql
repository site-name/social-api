CREATE TABLE IF NOT EXISTS checkout_lines (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  checkout_id varchar(36) NOT NULL,
  variant_id varchar(36) NOT NULL,
  quantity integer NOT NULL,
  is_gift boolean NOT NULL DEFAULT false,
  price_override decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  total_price_net_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  total_price_gross_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  tax_rate decimal(5,4)
);

CREATE INDEX idx_checkout_lines_checkout_id ON checkout_lines USING btree (checkout_id);
CREATE INDEX idx_checkout_lines_variant_id ON checkout_lines USING btree (variant_id);

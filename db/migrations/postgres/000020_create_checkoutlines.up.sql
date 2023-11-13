CREATE TABLE IF NOT EXISTS checkout_lines (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  checkout_id uuid NOT NULL,
  variant_id uuid NOT NULL,
  quantity integer NOT NULL
);

CREATE INDEX idx_checkout_lines_checkout_id ON checkout_lines USING btree (checkout_id);

CREATE INDEX idx_checkout_lines_variant_id ON checkout_lines USING btree (variant_id);
CREATE TABLE IF NOT EXISTS order_discounts (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid,
  type order_discount_type NOT NULL,
  value_type discount_value_type NOT NULL,
  value decimal(12,3) NOT NULL DEFAULT 0.00,
  amount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  name varchar(255),
  translated_name varchar(255),
  reason text
);

CREATE INDEX idx_order_discounts_name ON order_discounts USING btree (name);

CREATE INDEX idx_order_discounts_translated_name ON order_discounts USING btree (translated_name);

CREATE INDEX idx_order_discounts_name_lower_textpattern ON order_discounts USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_order_discounts_translated_name_lower_textpattern ON order_discounts USING btree (lower((translated_name)::text) text_pattern_ops);
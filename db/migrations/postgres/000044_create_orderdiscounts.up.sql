CREATE TABLE IF NOT EXISTS order_discounts (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid,
  type character varying(10) NOT NULL,
  value_type character varying(10) NOT NULL,
  value decimal(12,3) NOT NULL DEFAULT 0.00,
  amount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  currency varchar(3) NOT NULL,
  name character varying(255),
  translated_name character varying(255),
  reason text
);

CREATE INDEX idx_order_discounts_name ON order_discounts USING btree (name);

CREATE INDEX idx_order_discounts_translated_name ON order_discounts USING btree (translated_name);

CREATE INDEX idx_order_discounts_name_lower_textpattern ON order_discounts USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_order_discounts_translated_name_lower_textpattern ON order_discounts USING btree (lower((translated_name)::text) text_pattern_ops);
CREATE TABLE IF NOT EXISTS order_discounts (
  id varchar(36) NOT NULL PRIMARY KEY,
  order_id varchar(36),
  type order_discount_type NOT NULL,
  value_type discount_value_type NOT NULL,
  value decimal(12,3) NOT NULL DEFAULT 0.00,
  amount_value decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL,
  name varchar(255),
  translated_name varchar(255),
  reason text,
  promotion_rule_id varchar(36),
  voucher_id varchar(36),
  voucher_code varchar(255),
  created_at bigint NOT NULL
);

CREATE INDEX idx_order_discounts_name ON order_discounts USING btree (name);
CREATE INDEX idx_order_discounts_voucher_code ON order_discounts USING btree (voucher_code);
CREATE INDEX idx_order_discounts_promotion_rule ON order_discounts USING btree (promotion_rule_id);

CREATE INDEX idx_order_discounts_translated_name ON order_discounts USING btree (translated_name);

CREATE INDEX idx_order_discounts_name_lower_textpattern ON order_discounts USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_order_discounts_translated_name_lower_textpattern ON order_discounts USING btree (lower((translated_name)::text) text_pattern_ops);

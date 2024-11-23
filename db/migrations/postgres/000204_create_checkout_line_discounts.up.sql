CREATE TABLE IF NOT EXISTS checkout_line_discounts (
  id varchar(36) NOT NULL PRIMARY KEY,
  checkout_line_id varchar(36),
  type discount_type NOT NULL,
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
  created_at bigint NOT NULL,
  unique_type discount_type
);

ALTER TABLE checkout_line_discounts ADD CONSTRAINT fk_checkout_line_id FOREIGN KEY (checkout_line_id) REFERENCES checkout_lines(id) ON DELETE CASCADE;
ALTER TABLE checkout_line_discounts ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE SET NULL;
ALTER TABLE checkout_line_discounts ADD CONSTRAINT fk_voucher_id FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_checkout_line_discounts_voucher_code ON checkout_line_discounts USING btree (voucher_code);
CREATE INDEX IF NOT EXISTS idx_checkout_line_discounts_promotion_rule_id ON checkout_line_discounts USING BTREE (name);
CREATE UNIQUE INDEX idx_unique_checkout_line_id_unique_type ON checkout_line_discounts (checkout_line_id, unique_type);

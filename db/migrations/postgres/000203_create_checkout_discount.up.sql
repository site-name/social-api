CREATE TABLE IF NOT EXISTS checkout_discounts (
  id varchar(36) NOT NULL PRIMARY KEY,
  checkout_id varchar(36),
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
  created_at bigint NOT NULL
);

ALTER TABLE checkout_discounts ADD CONSTRAINT fk_checkout_id FOREIGN KEY (checkout_id) REFERENCES checkouts(token) ON DELETE CASCADE;
ALTER TABLE checkout_discounts ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE SET NULL;
ALTER TABLE checkout_discounts ADD CONSTRAINT fk_voucher_id FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_checkout_discounts_voucher_code ON checkout_discounts USING btree (voucher_code);
CREATE INDEX IF NOT EXISTS idx_checkout_discounts_name ON checkout_discounts USING btree (name);
CREATE INDEX IF NOT EXISTS idx_checkout_discounts_promotion_rule ON checkout_discounts USING btree (promotion_rule_id);

CREATE TABLE IF NOT EXISTS promotions (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(255) NOT NULL,
  type promotion_type NOT NULL default 'catalogue',
  description JSONB,
  start_date bigint NOT NULL, -- future time in milli
  "end_date" bigint,
  created_at bigint NOT NULL,
  updated_at bigint,
  last_modification_scheduled_at bigint,
  metadata jsonb,
  private_metadata jsonb
);

CREATE INDEX IF NOT EXISTS idx_promotions_start_date ON promotions USING btree (start_date);
CREATE INDEX IF NOT EXISTS idx_promotions_end_date ON promotions USING btree ("end_date");

--
CREATE TABLE IF NOT EXISTS promotion_rules (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(255) NOT NULL,
  description JSONB,
  promotion_id varchar(36) NOT NULL,
  catalogue_predicate JSONB,
  order_predicate JSONB,
  reward_value_type reward_value_type,
  reward_value decimal(12, 2),
  reward_type reward_type,
  variants_dirty boolean
);

ALTER TABLE promotion_rules ADD CONSTRAINT fk_promotion_id FOREIGN KEY (promotion_id) REFERENCES promotions(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS promotion_rule_channels (
  id varchar(36) NOT NULL PRIMARY KEY,
  promotion_rule_id varchar(36) NOT NULL,
  channel_id varchar(36) NOT NULL
);

ALTER TABLE promotion_rule_channels ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE CASCADE;
ALTER TABLE promotion_rule_channels ADD CONSTRAINT fk_channel_id FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS promotion_rule_product_variants (
  id varchar(36) NOT NULL PRIMARY KEY,
  promotion_rule_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL
);

ALTER TABLE promotion_rule_product_variants ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE CASCADE;
ALTER TABLE promotion_rule_product_variants ADD CONSTRAINT fk_product_variant_id FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS promotion_rule_gifts (
  id varchar(36) NOT NULL PRIMARY KEY,
  promotion_rule_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL
);

ALTER TABLE promotion_rule_gifts ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE CASCADE;
ALTER TABLE promotion_rule_gifts ADD CONSTRAINT fk_product_variant_id FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;

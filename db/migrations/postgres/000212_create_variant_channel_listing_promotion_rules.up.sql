CREATE TABLE IF NOT EXISTS variant_channel_listing_promotion_rules (
  id varchar(36) NOT NULL PRIMARY KEY,
  variant_channel_listing_id varchar(36) NOT NULL,
  promotion_rule_id varchar(36) NOT NULL,
  discount_amount decimal(12,3) NOT NULL DEFAULT 0.00,
  currency Currency NOT NULL
);

ALTER TABLE ONLY variant_channel_listing_promotion_rules
    ADD CONSTRAINT variant_channel_listing_promotion_rules_variant_channel_listing_id_promotion_rule_id_key UNIQUE (variant_channel_listing_id, promotion_rule_id);
ALTER TABLE ONLY variant_channel_listing_promotion_rules
    ADD CONSTRAINT fk_variant_channel_listing_id FOREIGN KEY (variant_channel_listing_id) REFERENCES product_variant_channel_listings(id) ON DELETE CASCADE;
ALTER TABLE ONLY variant_channel_listing_promotion_rules
    ADD CONSTRAINT fk_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE CASCADE;

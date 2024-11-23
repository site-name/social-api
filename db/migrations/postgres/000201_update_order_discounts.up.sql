ALTER TABLE ONLY order_discounts 
    ADD CONSTRAINT fk_order_discounts_promotion_rule_id FOREIGN KEY (promotion_rule_id) REFERENCES promotion_rules(id) ON DELETE CASCADE;

ALTER TABLE ONLY order_discounts 
    ADD CONSTRAINT fk_order_discounts_promotion_voucher_id FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE;


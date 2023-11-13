ALTER TABLE ONLY shipping_method_postal_code_rules
    ADD CONSTRAINT fk_shipping_method_postal_code_rules_shipping_methods FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id) ON DELETE CASCADE;

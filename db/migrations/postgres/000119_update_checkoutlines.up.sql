ALTER TABLE ONLY checkout_lines
    ADD CONSTRAINT fk_checkout_lines_checkouts FOREIGN KEY (checkout_id) REFERENCES checkouts(token) ON DELETE CASCADE;
ALTER TABLE ONLY checkout_lines
    ADD CONSTRAINT fk_checkout_lines_product_variants FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON DELETE CASCADE;
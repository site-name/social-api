ALTER TABLE ONLY checkout_lines
    ADD CONSTRAINT fk_checkout_lines_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token) ON DELETE CASCADE;
ALTER TABLE ONLY checkout_lines
    ADD CONSTRAINT fk_checkout_lines_product_variants FOREIGN KEY (variantid) REFERENCES product_variants(id) ON DELETE CASCADE;

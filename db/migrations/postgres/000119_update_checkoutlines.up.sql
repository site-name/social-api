ALTER TABLE ONLY checkoutlines
    ADD CONSTRAINT fk_checkoutlines_checkouts FOREIGN KEY (checkoutid) REFERENCES checkouts(token) ON DELETE CASCADE;
ALTER TABLE ONLY checkoutlines
    ADD CONSTRAINT fk_checkoutlines_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;

ALTER TABLE ONLY digital_contents
    ADD CONSTRAINT fk_digital_contents_product_variants FOREIGN KEY (productvariantid) REFERENCES product_variants(id) ON DELETE CASCADE;

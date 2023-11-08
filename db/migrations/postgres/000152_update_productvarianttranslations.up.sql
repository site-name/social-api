ALTER TABLE ONLY product_variant_translations
    ADD CONSTRAINT fk_product_variant_translations_product_variants FOREIGN KEY (productvariantid) REFERENCES product_variants(id) ON DELETE CASCADE;

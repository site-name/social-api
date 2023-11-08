ALTER TABLE ONLY variant_media
    ADD CONSTRAINT fk_variant_media_product_media FOREIGN KEY (mediaid) REFERENCES product_media(id) ON DELETE CASCADE;
ALTER TABLE ONLY variant_media
    ADD CONSTRAINT fk_variant_media_product_variants FOREIGN KEY (variantid) REFERENCES product_variants(id) ON DELETE CASCADE;

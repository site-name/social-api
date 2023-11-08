ALTER TABLE ONLY attribute_variants
    ADD CONSTRAINT fk_attribute_variants_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY attribute_variants
    ADD CONSTRAINT fk_attribute_variants_product_types FOREIGN KEY (producttypeid) REFERENCES product_types(id) ON DELETE CASCADE;

ALTER TABLE ONLY assigned_variant_attributes
    ADD CONSTRAINT fk_assigned_variant_attributes_attribute_variants FOREIGN KEY (assignmentid) REFERENCES attribute_variants(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_variant_attributes
    ADD CONSTRAINT fk_assigned_variant_attributes_product_variants FOREIGN KEY (variantid) REFERENCES product_variants(id) ON DELETE CASCADE;

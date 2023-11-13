ALTER TABLE ONLY assigned_variant_attribute_values
    ADD CONSTRAINT fk_assigned_variant_attribute_values_assigned_variant_attributes FOREIGN KEY (assignment_id) REFERENCES assigned_variant_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_variant_attribute_values
    ADD CONSTRAINT fk_assigned_variant_attribute_values_attribute_values FOREIGN KEY (value_id) REFERENCES attribute_values(id) ON DELETE CASCADE;
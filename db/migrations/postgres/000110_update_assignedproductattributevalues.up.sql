ALTER TABLE ONLY assigned_product_attribute_values
    ADD CONSTRAINT fk_assigned_product_attribute_values_assigned_product_attributes FOREIGN KEY (assignment_id) REFERENCES assigned_product_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_product_attribute_values
    ADD CONSTRAINT fk_assigned_product_attribute_values_attribute_values FOREIGN KEY (value_id) REFERENCES attribute_values(id) ON DELETE CASCADE;
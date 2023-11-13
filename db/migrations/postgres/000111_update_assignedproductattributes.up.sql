ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT fk_assigned_product_attributes_attribute_products FOREIGN KEY (assignment_id) REFERENCES attribute_products(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT fk_assigned_product_attributes_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
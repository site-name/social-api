ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT fk_assigned_product_attributes_category_attributes FOREIGN KEY (assignment_id) REFERENCES category_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_product_attributes
    ADD CONSTRAINT fk_assigned_product_attributes_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
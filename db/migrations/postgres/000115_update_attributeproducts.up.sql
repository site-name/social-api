ALTER TABLE ONLY attribute_products
    ADD CONSTRAINT fk_attribute_products_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY attribute_products
    ADD CONSTRAINT fk_attribute_products_product_types FOREIGN KEY (producttypeid) REFERENCES product_types(id) ON DELETE CASCADE;

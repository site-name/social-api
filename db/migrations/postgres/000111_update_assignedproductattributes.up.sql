ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_attributeproducts FOREIGN KEY (assignmentid) REFERENCES attributeproducts(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedproductattributes
    ADD CONSTRAINT fk_assignedproductattributes_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;

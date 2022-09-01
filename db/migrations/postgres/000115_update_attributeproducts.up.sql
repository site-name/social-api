ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT fk_attributeproducts_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT fk_attributeproducts_producttypes FOREIGN KEY (producttypeid) REFERENCES producttypes(id) ON DELETE CASCADE;

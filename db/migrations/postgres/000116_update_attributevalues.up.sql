ALTER TABLE ONLY attribute_values
    ADD CONSTRAINT fk_attribute_values_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;

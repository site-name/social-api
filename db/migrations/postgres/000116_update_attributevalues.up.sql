ALTER TABLE ONLY attribute_values
    ADD CONSTRAINT fk_attribute_values_attributes FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_page_attribute_values
    ADD CONSTRAINT fk_assigned_page_attribute_values_assigned_page_attributes FOREIGN KEY (assignmentid) REFERENCES assigned_page_attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_page_attribute_values
    ADD CONSTRAINT fk_assigned_page_attribute_values_attribute_values FOREIGN KEY (valueid) REFERENCES attribute_values(id) ON DELETE CASCADE;

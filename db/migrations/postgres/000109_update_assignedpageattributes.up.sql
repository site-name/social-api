ALTER TABLE ONLY assigned_page_attributes
    ADD CONSTRAINT fk_assigned_page_attributes_attribute_pages FOREIGN KEY (assignmentid) REFERENCES attribute_pages(id) ON DELETE CASCADE;
ALTER TABLE ONLY assigned_page_attributes
    ADD CONSTRAINT fk_assigned_page_attributes_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;

ALTER TABLE ONLY attribute_pages
    ADD CONSTRAINT fk_attribute_pages_page_types FOREIGN KEY (pagetypeid) REFERENCES page_types(id) ON DELETE CASCADE;

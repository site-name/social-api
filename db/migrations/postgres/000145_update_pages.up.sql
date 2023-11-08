ALTER TABLE ONLY pages
    ADD CONSTRAINT fk_pages_page_types FOREIGN KEY (pagetypeid) REFERENCES page_types(id) ON DELETE CASCADE;

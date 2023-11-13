ALTER TABLE ONLY pages
    ADD CONSTRAINT fk_pages_page_types FOREIGN KEY (page_type_id) REFERENCES page_types(id) ON DELETE CASCADE;
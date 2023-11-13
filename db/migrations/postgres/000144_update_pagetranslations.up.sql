ALTER TABLE ONLY page_translations
    ADD CONSTRAINT fk_page_translations_pages FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE;

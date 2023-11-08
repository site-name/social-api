ALTER TABLE ONLY page_translations
    ADD CONSTRAINT fk_page_translations_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;

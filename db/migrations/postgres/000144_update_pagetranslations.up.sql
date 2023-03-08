ALTER TABLE ONLY pagetranslations
    ADD CONSTRAINT fk_pagetranslations_pages FOREIGN KEY (pageid) REFERENCES pages(id) ON DELETE CASCADE;

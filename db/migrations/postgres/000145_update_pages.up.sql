ALTER TABLE ONLY pages
    ADD CONSTRAINT fk_pages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES pagetypes(id) ON DELETE CASCADE;

ALTER TABLE ONLY attributepages
    ADD CONSTRAINT fk_attributepages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES pagetypes(id) ON DELETE CASCADE;

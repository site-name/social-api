ALTER TABLE ONLY attributevariants
    ADD CONSTRAINT fk_attributevariants_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY attributevariants
    ADD CONSTRAINT fk_attributevariants_producttypes FOREIGN KEY (producttypeid) REFERENCES producttypes(id) ON DELETE CASCADE;

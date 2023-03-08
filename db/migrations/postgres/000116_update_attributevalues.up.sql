ALTER TABLE ONLY attributevalues
    ADD CONSTRAINT fk_attributevalues_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_attributevariants FOREIGN KEY (assignmentid) REFERENCES attributevariants(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;

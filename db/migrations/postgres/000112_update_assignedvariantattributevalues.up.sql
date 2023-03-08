ALTER TABLE ONLY assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_assignedvariantattributes FOREIGN KEY (assignmentid) REFERENCES assignedvariantattributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;

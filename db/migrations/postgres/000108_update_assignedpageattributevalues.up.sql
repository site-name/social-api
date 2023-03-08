ALTER TABLE ONLY assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_assignedpageattributes FOREIGN KEY (assignmentid) REFERENCES assignedpageattributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;

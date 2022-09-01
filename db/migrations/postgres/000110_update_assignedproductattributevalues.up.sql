ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_assignedproductattributes FOREIGN KEY (assignmentid) REFERENCES assignedproductattributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;

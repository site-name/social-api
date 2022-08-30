CREATE TABLE IF NOT EXISTS assignedproductattributevalues (
  id character varying(36) NOT NULL PRIMARY KEY,
  valueid character varying(36),
  assignmentid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT assignedproductattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);

ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_assignedproductattributes FOREIGN KEY (assignmentid) REFERENCES assignedproductattributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT fk_assignedproductattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;


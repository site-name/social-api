CREATE TABLE IF NOT EXISTS assignedvariantattributevalues (
  id character varying(36) NOT NULL PRIMARY KEY,
  valueid character varying(36),
  assignmentid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY assignedvariantattributevalues
    ADD CONSTRAINT assignedvariantattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);

ALTER TABLE ONLY assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_assignedvariantattributes FOREIGN KEY (assignmentid) REFERENCES assignedvariantattributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY assignedvariantattributevalues
    ADD CONSTRAINT fk_assignedvariantattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;


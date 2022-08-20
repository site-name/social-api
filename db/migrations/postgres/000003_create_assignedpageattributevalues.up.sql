CREATE TABLE IF NOT EXISTS assignedpageattributevalues (
  id character varying(36) NOT NULL PRIMARY KEY,
  valueid character varying(36),
  assignmentid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY assignedpageattributevalues
    ADD CONSTRAINT assignedpageattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);

ALTER TABLE ONLY assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_assignedpageattributes FOREIGN KEY (assignmentid) REFERENCES assignedpageattributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY assignedpageattributevalues
    ADD CONSTRAINT fk_assignedpageattributevalues_attributevalues FOREIGN KEY (valueid) REFERENCES attributevalues(id) ON DELETE CASCADE;



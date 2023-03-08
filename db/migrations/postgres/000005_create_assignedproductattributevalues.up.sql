CREATE TABLE IF NOT EXISTS assignedproductattributevalues (
  id character varying(36) NOT NULL PRIMARY KEY,
  valueid character varying(36),
  assignmentid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY assignedproductattributevalues
    ADD CONSTRAINT assignedproductattributevalues_valueid_assignmentid_key UNIQUE (valueid, assignmentid);


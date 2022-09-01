CREATE TABLE IF NOT EXISTS assignedvariantattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT assignedvariantattributes_variantid_assignmentid_key UNIQUE (variantid, assignmentid);


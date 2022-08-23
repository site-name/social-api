CREATE TABLE IF NOT EXISTS assignedvariantattributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT assignedvariantattributes_variantid_assignmentid_key UNIQUE (variantid, assignmentid);

ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_attributevariants FOREIGN KEY (assignmentid) REFERENCES attributevariants(id) ON DELETE CASCADE;

ALTER TABLE ONLY assignedvariantattributes
    ADD CONSTRAINT fk_assignedvariantattributes_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;


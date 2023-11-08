CREATE TABLE IF NOT EXISTS assigned_variant_attributes (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  assignmentid character varying(36)
);

ALTER TABLE ONLY assigned_variant_attributes
    ADD CONSTRAINT assigned_variant_attributes_variantid_assignmentid_key UNIQUE (variantid, assignmentid);


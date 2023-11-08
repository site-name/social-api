CREATE TABLE IF NOT EXISTS assigned_variant_attribute_values (
  id character varying(36) NOT NULL PRIMARY KEY,
  valueid character varying(36),
  assignmentid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY assigned_variant_attribute_values
    ADD CONSTRAINT assigned_variant_attribute_values_valueid_assignmentid_key UNIQUE (valueid, assignmentid);


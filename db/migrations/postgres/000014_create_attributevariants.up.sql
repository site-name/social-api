CREATE TABLE IF NOT EXISTS attributevariants (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  producttypeid character varying(36),
  variantselection boolean,
  sortorder integer
);

ALTER TABLE ONLY attributevariants
    ADD CONSTRAINT attributevariants_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);

ALTER TABLE ONLY attributevariants
    ADD CONSTRAINT fk_attributevariants_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY attributevariants
    ADD CONSTRAINT fk_attributevariants_producttypes FOREIGN KEY (producttypeid) REFERENCES producttypes(id) ON DELETE CASCADE;


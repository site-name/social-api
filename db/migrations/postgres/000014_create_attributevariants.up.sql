CREATE TABLE IF NOT EXISTS attribute_variants (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  producttypeid character varying(36),
  variantselection boolean,
  sortorder integer
);

ALTER TABLE ONLY attribute_variants
    ADD CONSTRAINT attribute_variants_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);


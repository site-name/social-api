CREATE TABLE IF NOT EXISTS attributeproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  producttypeid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT attributeproducts_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);


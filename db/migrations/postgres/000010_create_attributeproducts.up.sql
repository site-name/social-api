CREATE TABLE IF NOT EXISTS attributeproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  producttypeid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT attributeproducts_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);

ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT fk_attributeproducts_attributes FOREIGN KEY (attributeid) REFERENCES attributes(id) ON DELETE CASCADE;
ALTER TABLE ONLY attributeproducts
    ADD CONSTRAINT fk_attributeproducts_producttypes FOREIGN KEY (producttypeid) REFERENCES producttypes(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS attribute_products (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  producttypeid character varying(36),
  sortorder integer
);

ALTER TABLE ONLY attribute_products
    ADD CONSTRAINT attribute_products_attributeid_producttypeid_key UNIQUE (attributeid, producttypeid);


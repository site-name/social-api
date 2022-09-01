CREATE TABLE IF NOT EXISTS saleproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  productid character varying(36),
  createat bigint
);

ALTER TABLE ONLY saleproducts
    ADD CONSTRAINT saleproducts_saleid_productid_key UNIQUE (saleid, productid);

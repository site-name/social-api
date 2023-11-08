CREATE TABLE IF NOT EXISTS sale_products (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  productid character varying(36),
  createat bigint
);

ALTER TABLE ONLY sale_products
    ADD CONSTRAINT sale_products_saleid_productid_key UNIQUE (saleid, productid);

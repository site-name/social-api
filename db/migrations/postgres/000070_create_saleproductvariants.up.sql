CREATE TABLE IF NOT EXISTS saleproduct_variants (
  id character varying(36) NOT NULL PRIMARY KEY,
  saleid character varying(36),
  productvariantid character varying(36),
  createat bigint
);

ALTER TABLE ONLY saleproduct_variants
    ADD CONSTRAINT saleproduct_variants_saleid_productvariantid_key UNIQUE (saleid, productvariantid);

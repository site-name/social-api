CREATE TABLE IF NOT EXISTS voucherproduct_variants (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  productvariantid character varying(36),
  createat bigint
);

ALTER TABLE ONLY voucherproduct_variants
    ADD CONSTRAINT voucherproduct_variants_voucherid_productvariantid_key UNIQUE (voucherid, productvariantid);

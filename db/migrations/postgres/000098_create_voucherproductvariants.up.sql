CREATE TABLE IF NOT EXISTS voucherproductvariants (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  productvariantid character varying(36),
  createat bigint
);

ALTER TABLE ONLY voucherproductvariants
    ADD CONSTRAINT voucherproductvariants_voucherid_productvariantid_key UNIQUE (voucherid, productvariantid);

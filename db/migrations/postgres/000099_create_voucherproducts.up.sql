CREATE TABLE IF NOT EXISTS voucherproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT voucherproducts_voucherid_productid_key UNIQUE (voucherid, productid);

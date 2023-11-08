CREATE TABLE IF NOT EXISTS voucher_products (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT voucher_products_voucherid_productid_key UNIQUE (voucherid, productid);

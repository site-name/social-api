CREATE TABLE IF NOT EXISTS voucherproducts (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  productid character varying(36)
);

ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT voucherproducts_voucherid_productid_key UNIQUE (voucherid, productid);
ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT fk_voucherproducts_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE ONLY voucherproducts
    ADD CONSTRAINT fk_voucherproducts_vouchers FOREIGN KEY (voucherid) REFERENCES vouchers(id) ON DELETE CASCADE;

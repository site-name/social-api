CREATE TABLE IF NOT EXISTS voucher_customers (
  id character varying(36) NOT NULL PRIMARY KEY,
  voucherid character varying(36),
  customeremail character varying(128)
);

ALTER TABLE ONLY voucher_customers
    ADD CONSTRAINT voucher_customers_voucherid_customeremail_key UNIQUE (voucherid, customeremail);

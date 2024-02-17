CREATE TABLE IF NOT EXISTS voucher_products (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  product_id varchar(36) NOT NULL
);

ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT voucher_products_voucher_id_product_id_key UNIQUE (voucher_id, product_id);
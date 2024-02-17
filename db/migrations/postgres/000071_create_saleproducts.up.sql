CREATE TABLE IF NOT EXISTS sale_products (
  id varchar(36) NOT NULL PRIMARY KEY,
  sale_id varchar(36) NOT NULL,
  product_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_products
    ADD CONSTRAINT sale_products_sale_id_product_id_key UNIQUE (sale_id, product_id);
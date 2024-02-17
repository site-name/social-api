CREATE TABLE IF NOT EXISTS shipping_method_excluded_products (
  id varchar(36) NOT NULL PRIMARY KEY,
  shipping_method_id varchar(36) NOT NULL,
  product_id varchar(36) NOT NULL
);

ALTER TABLE ONLY shipping_method_excluded_products
    ADD CONSTRAINT shipping_method_excluded_products_shipping_method_id_product_id_key UNIQUE (shipping_method_id, product_id);
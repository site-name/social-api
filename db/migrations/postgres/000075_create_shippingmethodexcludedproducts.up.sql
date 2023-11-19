CREATE TABLE IF NOT EXISTS shipping_method_excluded_products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid NOT NULL,
  product_id uuid NOT NULL
);

ALTER TABLE ONLY shipping_method_excluded_products
    ADD CONSTRAINT shipping_method_excluded_products_shipping_method_id_product_id_key UNIQUE (shipping_method_id, product_id);
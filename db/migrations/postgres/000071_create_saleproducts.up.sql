CREATE TABLE IF NOT EXISTS sale_products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid NOT NULL,
  product_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_products
    ADD CONSTRAINT sale_products_sale_id_product_id_key UNIQUE (sale_id, product_id);
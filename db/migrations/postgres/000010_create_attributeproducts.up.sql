CREATE TABLE IF NOT EXISTS attribute_products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  attribute_id uuid NOT NULL,
  product_type_id uuid NOT NULL,
  sort_order integer
);

ALTER TABLE ONLY attribute_products
    ADD CONSTRAINT attribute_products_attribute_id_product_type_id_key UNIQUE (attribute_id, product_type_id);

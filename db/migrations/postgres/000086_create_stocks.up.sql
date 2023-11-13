CREATE TABLE IF NOT EXISTS stocks (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  warehouse_id uuid,
  product_variant_id uuid,
  quantity integer
);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT stocks_warehouse_id_product_variant_id_key UNIQUE (warehouse_id, product_variant_id);
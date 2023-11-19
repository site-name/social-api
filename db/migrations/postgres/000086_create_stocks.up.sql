CREATE TABLE IF NOT EXISTS stocks (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  warehouse_id uuid NOT NULL,
  product_variant_id uuid NOT NULL,
  quantity integer NOT NULL
);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT stocks_warehouse_id_product_variant_id_key UNIQUE (warehouse_id, product_variant_id);
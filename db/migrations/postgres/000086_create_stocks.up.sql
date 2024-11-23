CREATE TABLE IF NOT EXISTS stocks (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  warehouse_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL,
  quantity integer NOT NULL,
  quantity_allocated integer NOT NULL,
  annotations jsonb -- this is a JSON object that can store any additional data you want
);

ALTER TABLE ONLY stocks
    ADD CONSTRAINT stocks_warehouse_id_product_variant_id_key UNIQUE (warehouse_id, product_variant_id);
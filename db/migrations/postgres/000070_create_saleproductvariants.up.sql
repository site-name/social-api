CREATE TABLE IF NOT EXISTS sale_product_variants (
  id varchar(36) NOT NULL PRIMARY KEY,
  sale_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_product_variants
    ADD CONSTRAINT sale_product_variants_sale_id_product_variant_id_key UNIQUE (sale_id, product_variant_id);
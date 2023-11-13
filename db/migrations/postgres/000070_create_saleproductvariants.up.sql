CREATE TABLE IF NOT EXISTS sale_product_variants (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid,
  product_variant_id uuid,
  created_at bigint
);

ALTER TABLE ONLY sale_product_variants
    ADD CONSTRAINT sale_product_variants_sale_id_product_variant_id_key UNIQUE (sale_id, product_variant_id);
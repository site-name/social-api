CREATE TABLE IF NOT EXISTS voucherproduct_variants (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  voucher_id uuid NOT NULL,
  product_variant_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY voucherproduct_variants
    ADD CONSTRAINT voucherproduct_variants_voucher_id_product_variant_id_key UNIQUE (voucher_id, product_variant_id);
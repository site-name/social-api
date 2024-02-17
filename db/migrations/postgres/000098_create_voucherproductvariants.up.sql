CREATE TABLE IF NOT EXISTS voucherproduct_variants (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  product_variant_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY voucherproduct_variants
    ADD CONSTRAINT voucherproduct_variants_voucher_id_product_variant_id_key UNIQUE (voucher_id, product_variant_id);
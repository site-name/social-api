CREATE TABLE IF NOT EXISTS voucher_products (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  voucher_id uuid,
  product_id character varying(36)
);

ALTER TABLE ONLY voucher_products
    ADD CONSTRAINT voucher_products_voucher_id_product_id_key UNIQUE (voucher_id, product_id);
CREATE TABLE IF NOT EXISTS voucher_categories (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  voucher_id uuid,
  category_id uuid,
  created_at bigint
);

ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT voucher_categories_voucher_id_category_id_key UNIQUE (voucher_id, category_id);
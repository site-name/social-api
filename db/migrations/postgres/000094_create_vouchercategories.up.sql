CREATE TABLE IF NOT EXISTS voucher_categories (
  id varchar(36) NOT NULL PRIMARY KEY,
  voucher_id varchar(36) NOT NULL,
  category_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY voucher_categories
    ADD CONSTRAINT voucher_categories_voucher_id_category_id_key UNIQUE (voucher_id, category_id);
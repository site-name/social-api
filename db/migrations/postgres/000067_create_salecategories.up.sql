CREATE TABLE IF NOT EXISTS sale_categories (
  id varchar(36) NOT NULL PRIMARY KEY,
  sale_id varchar(36) NOT NULL,
  category_id varchar(36) NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT sale_categories_sale_id_category_id_key UNIQUE (sale_id, category_id);
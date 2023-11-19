CREATE TABLE IF NOT EXISTS sale_categories (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id uuid NOT NULL,
  category_id uuid NOT NULL,
  created_at bigint NOT NULL
);

ALTER TABLE ONLY sale_categories
    ADD CONSTRAINT sale_categories_sale_id_category_id_key UNIQUE (sale_id, category_id);
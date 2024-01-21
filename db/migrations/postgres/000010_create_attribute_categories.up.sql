CREATE TABLE IF NOT EXISTS category_attributes (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  attribute_id uuid NOT NULL,
  category_id uuid NOT NULL,
  sort_order integer
);

ALTER TABLE ONLY category_attributes
    ADD CONSTRAINT category_attributes_attribute_id_category_id_key UNIQUE (attribute_id, category_id);

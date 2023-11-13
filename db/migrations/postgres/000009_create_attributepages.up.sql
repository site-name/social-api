CREATE TABLE IF NOT EXISTS attribute_pages (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  attribute_id uuid NOT NULL,
  page_type_id uuid NOT NULL,
  sort_order integer
);

ALTER TABLE ONLY attribute_pages
    ADD CONSTRAINT attribute_pages_attribute_id_page_type_id_key UNIQUE (attribute_id, page_type_id);
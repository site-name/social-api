CREATE TABLE IF NOT EXISTS attribute_pages (
  id varchar(36) NOT NULL PRIMARY KEY,
  attribute_id varchar(36) NOT NULL,
  page_type_id varchar(36) NOT NULL,
  sort_order integer
);

ALTER TABLE ONLY attribute_pages
    ADD CONSTRAINT attribute_pages_attribute_id_page_type_id_key UNIQUE (attribute_id, page_type_id);
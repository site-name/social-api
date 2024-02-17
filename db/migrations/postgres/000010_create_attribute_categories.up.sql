CREATE TABLE IF NOT EXISTS category_attributes (
  id varchar(36) NOT NULL PRIMARY KEY,
  attribute_id varchar(36) NOT NULL,
  category_id varchar(36) NOT NULL,
  sort_order integer
);

ALTER TABLE ONLY category_attributes
    ADD CONSTRAINT category_attributes_attribute_id_category_id_key UNIQUE (attribute_id, category_id);

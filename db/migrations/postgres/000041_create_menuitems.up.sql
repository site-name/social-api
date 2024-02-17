CREATE TABLE IF NOT EXISTS menu_items (
  id varchar(36) NOT NULL PRIMARY KEY,
  menu_id varchar(36) NOT NULL,
  name varchar(128) NOT NULL,
  parent_id varchar(36),
  url varchar(256),
  category_id varchar(36),
  collection_id varchar(36),
  page_id varchar(36),
  metadata jsonb,
  private_metadata jsonb,
  sort_order integer
);
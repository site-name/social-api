CREATE TABLE IF NOT EXISTS menu_items (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  menu_id uuid,
  name character varying(128),
  parent_id uuid,
  url character varying(256),
  category_id uuid,
  collection_id uuid,
  page_id uuid,
  metadata jsonb,
  private_metadata jsonb,
  sort_order integer
);
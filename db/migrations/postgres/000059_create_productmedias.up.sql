CREATE TABLE IF NOT EXISTS product_media (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  product_id uuid NOT NULL,
  ppoi varchar(20) NOT NULL,
  image varchar(200) NOT NULL,
  alt varchar(128) NOT NULL,
  type varchar(32) NOT NULL,
  external_url varchar(256),
  oembed_data jsonb,
  sort_order integer
);
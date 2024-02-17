CREATE TABLE IF NOT EXISTS product_media (
  id varchar(36) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  product_id varchar(36) NOT NULL,
  ppoi varchar(20) NOT NULL,
  image varchar(200) NOT NULL,
  alt varchar(128) NOT NULL,
  type varchar(32) NOT NULL,
  external_url varchar(256),
  oembed_data jsonb,
  sort_order integer
);
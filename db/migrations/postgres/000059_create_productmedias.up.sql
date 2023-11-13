CREATE TABLE IF NOT EXISTS product_media (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint,
  product_id uuid,
  ppoi character varying(20),
  image character varying(200),
  alt character varying(128),
  type character varying(32),
  external_url character varying(256),
  oembed_data text,
  sort_order integer
);
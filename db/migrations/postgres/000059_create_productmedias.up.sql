CREATE TABLE IF NOT EXISTS product_media (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at bigint NOT NULL,
  product_id uuid NOT NULL,
  ppoi character varying(20) NOT NULL,
  image character varying(200) NOT NULL,
  alt character varying(128) NOT NULL,
  type character varying(32) NOT NULL,
  external_url character varying(256),
  oembed_data jsonb,
  sort_order integer
);
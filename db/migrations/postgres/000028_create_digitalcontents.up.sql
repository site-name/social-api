CREATE TABLE IF NOT EXISTS digital_contents (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  use_default_settings boolean,
  automatic_fulfillment boolean,
  content_type character varying(128),
  product_variant_id uuid,
  content_file character varying(200),
  max_downloads integer,
  url_valid_days integer,
  metadata jsonb,
  private_metadata jsonb
);
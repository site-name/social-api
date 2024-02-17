CREATE TABLE IF NOT EXISTS digital_contents (
  id varchar(36) NOT NULL PRIMARY KEY,
  use_default_settings boolean,
  automatic_fulfillment boolean,
  content_type content_type NOT NULL,
  product_variant_id varchar(36) NOT NULL,
  content_file varchar(200) NOT NULL,
  max_downloads integer,
  url_valid_days integer,
  metadata jsonb,
  private_metadata jsonb
);